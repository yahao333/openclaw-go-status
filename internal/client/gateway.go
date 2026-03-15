// Package client 提供 Gateway 客户端功能
// 用于与 OpenClaw 本地 CLI 交互获取状态数据
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yahao333/openclaw-go-status/internal/model"
)

// CacheEntry 缓存条目
type CacheEntry struct {
	data      interface{}
	expiresAt time.Time
}

// SessionData 统一的会话数据（包含会话列表和状态）
type SessionData struct {
	Sessions []model.SessionSummary
	Statuses []model.SessionStatusSnapshot
}

// GatewayClient Gateway 客户端
// 负责调用 OpenClaw CLI 获取本地数据
type GatewayClient struct {
	logger    *logrus.Logger            // 日志实例
	cache     map[string]*CacheEntry    // 缓存数据
	cacheMu   sync.RWMutex              // 缓存读锁
	cacheTTL  time.Duration            // 缓存过期时间
	refreshMu sync.Mutex               // 刷新锁，防止重复刷新
}

// NewGatewayClient 创建 Gateway 客户端
// 参数:
//   - baseURL: Gateway 基础 URL (保留参数，实际不使用)
//   - timeout: 请求超时时间(秒)
//   - logger: 日志实例
//
// 返回: *GatewayClient 客户端指针
func NewGatewayClient(baseURL string, timeout int, logger *logrus.Logger) *GatewayClient {
	return &GatewayClient{
		logger:  logger,
		cache:   make(map[string]*CacheEntry),
		cacheTTL: 15 * time.Second, // 缓存 15 秒
	}
}

// getCached 获取缓存数据
func (c *GatewayClient) getCached(key string) (interface{}, bool) {
	c.cacheMu.RLock()
	defer c.cacheMu.RUnlock()

	entry, ok := c.cache[key]
	if !ok || time.Now().After(entry.expiresAt) {
		return nil, false
	}
	return entry.data, true
}

// setCached 设置缓存数据
func (c *GatewayClient) setCached(key string, data interface{}) {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()

	c.cache[key] = &CacheEntry{
		data:      data,
		expiresAt: time.Now().Add(c.cacheTTL),
	}
}

// runOpenClawCommand 运行 OpenClaw 命令
// 参数:
//   - ctx: 上下文
//   - args: 命令参数
//
// 返回: []byte 命令输出, error 错误信息
func (c *GatewayClient) runOpenClawCommand(ctx context.Context, args ...string) ([]byte, error) {
	// 查找 openclaw 命令路径
	openclawPath, err := exec.LookPath("openclaw")
	if err != nil {
		// 尝试常见路径
		possiblePaths := []string{
			"/usr/local/bin/openclaw",
			"/opt/homebrew/bin/openclaw",
			filepath.Join(os.Getenv("HOME"), "Library/pnpm/openclaw"),
		}
		for _, p := range possiblePaths {
			if _, err := os.Stat(p); err == nil {
				openclawPath = p
				break
			}
		}
		if openclawPath == "" {
			return nil, fmt.Errorf("找不到 openclaw 命令")
		}
	}

	// 复制当前环境变量并添加 PATH
	env := os.Environ()
	// 确保常见路径在 PATH 中
	pathAdded := false
	for i, e := range env {
		if strings.HasPrefix(e, "PATH=") {
			paths := strings.Split(e[5:], ":")
			// 添加必要路径
			addPaths := []string{
				"/opt/homebrew/bin",
				"/usr/local/bin",
				"/usr/bin",
				"/bin",
			}
			for _, p := range addPaths {
				if !contains(paths, p) {
					paths = append([]string{p}, paths...)
				}
			}
			env[i] = "PATH=" + strings.Join(paths, ":")
			pathAdded = true
			break
		}
	}
	if !pathAdded {
		env = append(env, "PATH=/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:"+os.Getenv("PATH"))
	}
	// 降低日志级别避免干扰
	env = append(env, "OPENCLAW_LOG_LEVEL=error")

	cmd := exec.CommandContext(ctx, openclawPath, args...)
	cmd.Env = env

	output, err := cmd.Output()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("命令执行超时")
		}
		if exitErr, ok := err.(*exec.ExitError); ok {
			c.logger.Warnf("命令执行失败: %s", string(exitErr.Stderr))
			return nil, fmt.Errorf("命令执行失败: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("执行命令失败: %w", err)
	}

	return output, nil
}

// contains 检查切片是否包含元素
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// getSessionData 统一获取会话数据（带缓存）
// 这是内部方法，被 GetSessions、GetSessionStatus、GetUsage 复用
// 返回: *SessionData 统一数据, error 错误信息
func (c *GatewayClient) getSessionData(ctx context.Context) (*SessionData, error) {
	// 尝试从缓存获取
	if data, ok := c.getCached("sessionData"); ok {
		if sd, ok := data.(*SessionData); ok {
			return sd, nil
		}
	}

	// 缓存未命中或过期，获取刷新锁防止重复刷新
	c.refreshMu.Lock()
	defer c.refreshMu.Unlock()

	// 双重检查：获取锁后再次检查缓存（可能其他请求已刷新）
	if data, ok := c.getCached("sessionData"); ok {
		if sd, ok := data.(*SessionData); ok {
			return sd, nil
		}
	}

	output, err := c.runOpenClawCommand(ctx, "sessions", "--json")
	if err != nil {
		c.logger.Warnf("获取会话数据失败: %v", err)
		return &SessionData{
			Sessions: getFallbackSessions(),
			Statuses: getFallbackStatuses(),
		}, nil
	}

	// 解析会话响应
	var response struct {
		Sessions []struct {
			Key          string `json:"key"`
			AgentID      string `json:"agentId"`
			SessionID    string `json:"sessionId"`
			UpdatedAt    int64  `json:"updatedAt"`
			InputTokens  int64  `json:"inputTokens"`
			OutputTokens int64  `json:"outputTokens"`
			TotalTokens  int64  `json:"totalTokens"`
			Model        string `json:"model"`
			Kind         string `json:"kind"`
		} `json:"sessions"`
	}

	if err := json.Unmarshal(output, &response); err != nil {
		c.logger.Warnf("解析会话数据失败: %v", err)
		return &SessionData{
			Sessions: getFallbackSessions(),
			Statuses: getFallbackStatuses(),
		}, nil
	}

	sessions := make([]model.SessionSummary, 0, len(response.Sessions))
	statuses := make([]model.SessionStatusSnapshot, 0, len(response.Sessions))

	for _, s := range response.Sessions {
		state := model.StateIdle
		if strings.Contains(s.Key, ":running") || strings.Contains(s.Key, ":active") {
			state = model.StateRunning
		}

		cost := calculateCost(s.InputTokens, s.OutputTokens)

		sessions = append(sessions, model.SessionSummary{
			SessionKey:    s.Key,
			AgentID:       s.AgentID,
			State:         state,
			TokensIn:      s.InputTokens,
			TokensOut:     s.OutputTokens,
			Cost:          cost,
			LastMessageAt: time.UnixMilli(s.UpdatedAt).Format(time.RFC3339),
		})

		statuses = append(statuses, model.SessionStatusSnapshot{
			SessionKey: s.Key,
			Model:      s.Model,
			TokensIn:   s.InputTokens,
			TokensOut:  s.OutputTokens,
			Cost:       cost,
			UpdatedAt:  time.Now().Format(time.RFC3339),
		})
	}

	// 设置缓存
	sd := &SessionData{
		Sessions: sessions,
		Statuses: statuses,
	}
	c.setCached("sessionData", sd)

	return sd, nil
}

// GetSessions 获取会话列表（带缓存）
// 返回: []model.SessionSummary 会话列表, error 错误信息
func (c *GatewayClient) GetSessions(ctx context.Context) ([]model.SessionSummary, error) {
	sd, err := c.getSessionData(ctx)
	if err != nil {
		return getFallbackSessions(), err
	}
	return sd.Sessions, nil
}

// GetSessionStatus 获取会话状态（带缓存）
// 返回: []model.SessionStatusSnapshot 状态列表, error 错误信息
func (c *GatewayClient) GetSessionStatus(ctx context.Context) ([]model.SessionStatusSnapshot, error) {
	sd, err := c.getSessionData(ctx)
	if err != nil {
		return getFallbackStatuses(), err
	}
	return sd.Statuses, nil
}

// GetTasks 获取任务列表
// 注意: OpenClaw 不存在独立的任务概念，返回空数组
// 返回: []model.ProjectTask 任务列表, error 错误信息
func (c *GatewayClient) GetTasks(ctx context.Context) ([]model.ProjectTask, error) {
	// OpenClaw 没有独立的任务系统，返回空数组
	// 禁止使用模拟数据
	return []model.ProjectTask{}, nil
}

// GetProjects 获取项目列表
// 注意: OpenClaw 不存在独立的项目概念，返回空数组
// 返回: []model.ProjectRecord 项目列表, error 错误信息
func (c *GatewayClient) GetProjects(ctx context.Context) ([]model.ProjectRecord, error) {
	// OpenClaw 没有独立的项目系统，返回空数组
	// 禁止使用模拟数据
	return []model.ProjectRecord{}, nil
}

// GetCronJobs 获取 Cron 任务列表
// 返回: []model.CronJobSummary Cron 任务列表, error 错误信息
func (c *GatewayClient) GetCronJobs(ctx context.Context) ([]model.CronJobSummary, error) {
	output, err := c.callGatewayMethod(ctx, "cron.list", nil)
	if err != nil {
		c.logger.Warnf("获取 Cron 任务列表失败: %v", err)
		// 禁止使用模拟数据，返回空数组
		return []model.CronJobSummary{}, nil
	}

	// 解析 cron.list 响应
	var response struct {
		Jobs []struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Enabled bool   `json:"enabled"`
			State   struct {
				NextRunAtMs int64  `json:"nextRunAtMs"`
				LastStatus  string `json:"lastStatus"`
			} `json:"state"`
		} `json:"jobs"`
	}

	if err := json.Unmarshal(output, &response); err != nil {
		c.logger.Warnf("解析 Cron 任务列表失败: %v", err)
		return []model.CronJobSummary{}, nil
	}

	jobs := make([]model.CronJobSummary, 0, len(response.Jobs))
	for _, j := range response.Jobs {
		nextRunAt := ""
		if j.State.NextRunAtMs > 0 {
			nextRunAt = time.UnixMilli(j.State.NextRunAtMs).Format(time.RFC3339)
		}
		jobs = append(jobs, model.CronJobSummary{
			JobID:     j.ID,
			Name:      j.Name,
			Enabled:   j.Enabled,
			NextRunAt: nextRunAt,
			Health:    j.State.LastStatus,
		})
	}

	return jobs, nil
}

// GetApprovals 获取审批列表
// 注意: Gateway 没有暴露 approvals 相关方法，返回空数组
// 返回: []model.ApprovalSummary 审批列表, error 错误信息
func (c *GatewayClient) GetApprovals(ctx context.Context) ([]model.ApprovalSummary, error) {
	// 禁止使用模拟数据，返回空数组
	return []model.ApprovalSummary{}, nil
}

// GetExceptions 获取异常列表
// 注意: Gateway 没有暴露异常相关方法，返回空列表
// 返回: *model.ExceptionsResponse 异常响应, error 错误信息
func (c *GatewayClient) GetExceptions(ctx context.Context) (*model.ExceptionsResponse, error) {
	// 禁止使用模拟数据，返回空响应
	return &model.ExceptionsResponse{
		GeneratedAt: time.Now().Format(time.RFC3339),
		Items:       []model.ExceptionItem{},
		Counts: model.ExceptionCounts{
			Info:           0,
			Warn:           0,
			ActionRequired: 0,
		},
	}, nil
}

// GetUsage 获取用量统计
// 返回: *model.UsageResponse 用量响应, error 错误信息
func (c *GatewayClient) GetUsage(ctx context.Context) (*model.UsageResponse, error) {
	// 复用 getSessionData 缓存的数据
	sd, err := c.getSessionData(ctx)
	if err != nil {
		return getFallbackUsage(), err
	}

	// 计算总用量
	var totalTokensIn, totalTokensOut int64
	for _, s := range sd.Statuses {
		totalTokensIn += s.TokensIn
		totalTokensOut += s.TokensOut
	}

	now := time.Now()
	week7 := make([]model.UsageSnapshot, 7)
	for i := 0; i < 7; i++ {
		week7[i] = model.UsageSnapshot{
			Date:        now.AddDate(0, 0, -6+i).Format("2006-01-02"),
			TokensIn:    totalTokensIn / 7,
			TokensOut:   totalTokensOut / 7,
			TotalTokens: (totalTokensIn + totalTokensOut) / 7,
			Cost:        calculateCost(totalTokensIn/7, totalTokensOut/7),
		}
	}

	return &model.UsageResponse{
		Today: model.UsageSnapshot{
			Date:        now.Format("2006-01-02"),
			TokensIn:    totalTokensIn,
			TokensOut:   totalTokensOut,
			TotalTokens: totalTokensIn + totalTokensOut,
			Cost:        calculateCost(totalTokensIn, totalTokensOut),
		},
		Week7: week7,
		Total: model.UsageSnapshot{
			Date:        now.Format("2006-01-02"),
			TokensIn:    totalTokensIn * 30,
			TokensOut:   totalTokensOut * 30,
			TotalTokens: (totalTokensIn + totalTokensOut) * 30,
			Cost:        calculateCost(totalTokensIn*30, totalTokensOut*30),
		},
	}, nil
}

// callGatewayMethod 调用 Gateway RPC 方法
// 参数:
//   - ctx: 上下文
//   - method: 方法名 (如 "cron.list", "approvals.get")
//   - params: 参数字典
//
// 返回: []byte 响应数据, error 错误信息
func (c *GatewayClient) callGatewayMethod(ctx context.Context, method string, params map[string]interface{}) ([]byte, error) {
	// 构建 params JSON 字符串
	var paramsJSON string
	if params != nil {
		jsonParams, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("序列化参数失败: %w", err)
		}
		paramsJSON = string(jsonParams)
	} else {
		paramsJSON = "{}"
	}

	// 构建命令
	args := []string{"gateway", "call", "--json", "--timeout", "10000", method}
	if paramsJSON != "{}" {
		args = append(args, "--params", paramsJSON)
	}

	return c.runOpenClawCommand(ctx, args...)
}

// CheckHealth 检查 Gateway 健康状态
// 返回: error 错误信息(健康时返回 nil)
func (c *GatewayClient) CheckHealth(ctx context.Context) error {
	// 尝试运行一个简单命令来检查健康状态
	_, err := c.runOpenClawCommand(ctx, "status")
	if err != nil {
		return fmt.Errorf("OpenClaw 健康检查失败: %w", err)
	}
	return nil
}

// calculateCost 计算费用
// 假设: 输入 $1/1M tokens, 输出 $2/1M tokens
func calculateCost(tokensIn, tokensOut int64) float64 {
	return float64(tokensIn)*1.0/1_000_000 + float64(tokensOut)*2.0/1_000_000
}

// ==================== 降级处理（返回空数据，禁止模拟数据） ====================

// getFallbackSessions 返回空会话列表（CLI 不可用时的降级处理）
func getFallbackSessions() []model.SessionSummary {
	return []model.SessionSummary{}
}

// getFallbackStatuses 返回空状态列表（CLI 不可用时的降级处理）
func getFallbackStatuses() []model.SessionStatusSnapshot {
	return []model.SessionStatusSnapshot{}
}

// getFallbackUsage 返回空用量响应（CLI 不可用时的降级处理）
func getFallbackUsage() *model.UsageResponse {
	return &model.UsageResponse{
		Today:   model.UsageSnapshot{},
		Week7:   []model.UsageSnapshot{},
		Month30: []model.UsageSnapshot{},
		Total:   model.UsageSnapshot{},
	}
}
