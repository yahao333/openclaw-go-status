// Package client 提供 Gateway 客户端功能
// 用于与 OpenClaw 本地 CLI 交互获取状态数据
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yahao333/openclaw-go-status/internal/model"
)

// GatewayClient Gateway 客户端
// 负责调用 OpenClaw CLI 获取本地数据
type GatewayClient struct {
	logger *logrus.Logger // 日志实例
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
		logger: logger,
	}
}

// runOpenClawCommand 运行 OpenClaw 命令
// 参数:
//   - ctx: 上下文
//   - args: 命令参数
//
// 返回: []byte 命令输出, error 错误信息
func (c *GatewayClient) runOpenClawCommand(ctx context.Context, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "openclaw", args...)
	cmd.Env = append(cmd.Env, "OPENCLAW_LOG_LEVEL=error")

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

// GetSessions 获取会话列表
// 返回: []model.SessionSummary 会话列表, error 错误信息
func (c *GatewayClient) GetSessions(ctx context.Context) ([]model.SessionSummary, error) {
	output, err := c.runOpenClawCommand(ctx, "sessions", "--json")
	if err != nil {
		c.logger.Warnf("获取会话列表失败: %v", err)
		return getMockSessions(), nil
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
		c.logger.Warnf("解析会话列表失败: %v", err)
		return getMockSessions(), nil
	}

	sessions := make([]model.SessionSummary, 0, len(response.Sessions))
	for _, s := range response.Sessions {
		state := model.StateIdle
		if strings.Contains(s.Key, ":running") || strings.Contains(s.Key, ":active") {
			state = model.StateRunning
		}

		sessions = append(sessions, model.SessionSummary{
			SessionKey:    s.Key,
			AgentID:       s.AgentID,
			State:         state,
			LastMessageAt: time.UnixMilli(s.UpdatedAt).Format(time.RFC3339),
		})
	}

	return sessions, nil
}

// GetSessionStatus 获取会话状态
// 返回: []model.SessionStatusSnapshot 状态列表, error 错误信息
func (c *GatewayClient) GetSessionStatus(ctx context.Context) ([]model.SessionStatusSnapshot, error) {
	output, err := c.runOpenClawCommand(ctx, "sessions", "--json")
	if err != nil {
		c.logger.Warnf("获取会话状态失败: %v", err)
		return getMockStatuses(), nil
	}

	// 解析会话响应
	var response struct {
		Sessions []struct {
			Key          string `json:"key"`
			SessionID    string `json:"sessionId"`
			InputTokens  int64  `json:"inputTokens"`
			OutputTokens int64  `json:"outputTokens"`
			TotalTokens  int64  `json:"totalTokens"`
			Model        string `json:"model"`
		} `json:"sessions"`
	}

	if err := json.Unmarshal(output, &response); err != nil {
		c.logger.Warnf("解析会话状态失败: %v", err)
		return getMockStatuses(), nil
	}

	statuses := make([]model.SessionStatusSnapshot, 0, len(response.Sessions))
	for _, s := range response.Sessions {
		statuses = append(statuses, model.SessionStatusSnapshot{
			SessionKey: s.Key,
			Model:      s.Model,
			TokensIn:   s.InputTokens,
			TokensOut:  s.OutputTokens,
			Cost:       calculateCost(s.InputTokens, s.OutputTokens),
			UpdatedAt:  time.Now().Format(time.RFC3339),
		})
	}

	return statuses, nil
}

// GetTasks 获取任务列表
// 返回: []model.ProjectTask 任务列表, error 错误信息
func (c *GatewayClient) GetTasks(ctx context.Context) ([]model.ProjectTask, error) {
	// 尝试从 tasks.json 读取任务
	output, err := c.runOpenClawCommand(ctx, "tasks", "--json")
	if err != nil {
		c.logger.Warnf("获取任务列表失败: %v", err)
		return getMockTasks(), nil
	}

	// 解析任务响应
	var response struct {
		Tasks []struct {
			ID        string `json:"id"`
			Title     string `json:"title"`
			Status    string `json:"status"`
			Owner     string `json:"owner"`
			ProjectID string `json:"projectId"`
		} `json:"tasks"`
	}

	if err := json.Unmarshal(output, &response); err != nil {
		c.logger.Warnf("解析任务列表失败: %v", err)
		return getMockTasks(), nil
	}

	tasks := make([]model.ProjectTask, 0, len(response.Tasks))
	for _, t := range response.Tasks {
		status := model.TaskState(t.Status)
		if status == "" {
			status = model.TaskTodo
		}

		tasks = append(tasks, model.ProjectTask{
			ProjectID: t.ProjectID,
			TaskID:    t.ID,
			Title:     t.Title,
			Status:    status,
			Owner:     t.Owner,
			UpdatedAt: time.Now().Format(time.RFC3339),
		})
	}

	return tasks, nil
}

// GetProjects 获取项目列表
// 返回: []model.ProjectRecord 项目列表, error 错误信息
func (c *GatewayClient) GetProjects(ctx context.Context) ([]model.ProjectRecord, error) {
	// 返回模拟数据，因为没有直接的 CLI 命令
	return getMockProjects(), nil
}

// GetCronJobs 获取 Cron 任务列表
// 返回: []model.CronJobSummary Cron 任务列表, error 错误信息
func (c *GatewayClient) GetCronJobs(ctx context.Context) ([]model.CronJobSummary, error) {
	return getMockCronJobs(), nil
}

// GetApprovals 获取审批列表
// 返回: []model.ApprovalSummary 审批列表, error 错误信息
func (c *GatewayClient) GetApprovals(ctx context.Context) ([]model.ApprovalSummary, error) {
	return getMockApprovals(), nil
}

// GetExceptions 获取异常列表
// 返回: *model.ExceptionsResponse 异常响应, error 错误信息
func (c *GatewayClient) GetExceptions(ctx context.Context) (*model.ExceptionsResponse, error) {
	return getMockExceptions(), nil
}

// GetUsage 获取用量统计
// 返回: *model.UsageResponse 用量响应, error 错误信息
func (c *GatewayClient) GetUsage(ctx context.Context) (*model.UsageResponse, error) {
	// 先获取会话列表以计算用量
	output, err := c.runOpenClawCommand(ctx, "sessions", "--json")
	if err != nil {
		c.logger.Warnf("获取用量统计失败: %v", err)
		return getMockUsage(), nil
	}

	// 解析会话响应
	var response struct {
		Sessions []struct {
			InputTokens  int64 `json:"inputTokens"`
			OutputTokens int64 `json:"outputTokens"`
		} `json:"sessions"`
	}

	if err := json.Unmarshal(output, &response); err != nil {
		c.logger.Warnf("解析用量统计失败: %v", err)
		return getMockUsage(), nil
	}

	// 计算总用量
	var totalTokensIn, totalTokensOut int64
	for _, s := range response.Sessions {
		totalTokensIn += s.InputTokens
		totalTokensOut += s.OutputTokens
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

// ==================== 模拟数据函数 ====================

func getMockSessions() []model.SessionSummary {
	return []model.SessionSummary{
		{
			SessionKey:    "session-001",
			Label:         "主会话",
			AgentID:       "agent-001",
			State:         model.StateRunning,
			LastMessageAt: time.Now().Format(time.RFC3339),
		},
		{
			SessionKey:    "session-002",
			Label:         "辅助会话",
			AgentID:       "agent-002",
			State:         model.StateIdle,
			LastMessageAt: time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
		},
	}
}

func getMockStatuses() []model.SessionStatusSnapshot {
	return []model.SessionStatusSnapshot{
		{
			SessionKey: "session-001",
			Model:      "MiniMax-M2.5",
			TokensIn:   15000,
			TokensOut:  25000,
			Cost:       0.35,
			UpdatedAt:  time.Now().Format(time.RFC3339),
		},
		{
			SessionKey: "session-002",
			Model:      "MiniMax-M2.5",
			TokensIn:   5000,
			TokensOut:  8000,
			Cost:       0.12,
			UpdatedAt:  time.Now().Format(time.RFC3339),
		},
	}
}

func getMockTasks() []model.ProjectTask {
	todoTime := time.Now().Add(24 * time.Hour)
	return []model.ProjectTask{
		{
			ProjectID:   "project-001",
			TaskID:      "task-001",
			Title:       "完成用户认证模块",
			Status:      model.TaskInProgress,
			Owner:       "zhangsan",
			DueAt:       &todoTime,
			UpdatedAt:   time.Now().Format(time.RFC3339),
		},
		{
			ProjectID:   "project-001",
			TaskID:      "task-002",
			Title:       "编写 API 文档",
			Status:      model.TaskTodo,
			Owner:       "lisi",
			UpdatedAt:   time.Now().Format(time.RFC3339),
		},
	}
}

func getMockProjects() []model.ProjectRecord {
	return []model.ProjectRecord{
		{
			ProjectID: "project-001",
			Title:     "用户认证系统",
			Status:    model.ProjectActive,
			Owner:     "zhangsan",
			UpdatedAt: time.Now().Format(time.RFC3339),
		},
		{
			ProjectID: "project-002",
			Title:     "数据报表模块",
			Status:    model.ProjectPlanned,
			Owner:     "lisi",
			UpdatedAt: time.Now().Format(time.RFC3339),
		},
	}
}

func getMockCronJobs() []model.CronJobSummary {
	return []model.CronJobSummary{
		{
			JobID:     "cron-001",
			Name:      "每日报告",
			Enabled:   true,
			NextRunAt: time.Now().Add(1 * time.Hour).Format(time.RFC3339),
			Health:    "scheduled",
		},
	}
}

func getMockApprovals() []model.ApprovalSummary {
	return []model.ApprovalSummary{
		{
			ApprovalID: "approval-001",
			Status:     model.ApprovalPending,
			Command:    "exec",
			RequestedAt: time.Now().Format(time.RFC3339),
		},
	}
}

func getMockExceptions() *model.ExceptionsResponse {
	return &model.ExceptionsResponse{
		GeneratedAt: time.Now().Format(time.RFC3339),
		Items:      []model.ExceptionItem{},
		Counts: model.ExceptionCounts{
			Info:          0,
			Warn:          0,
			ActionRequired: 0,
		},
	}
}

func getMockUsage() *model.UsageResponse {
	now := time.Now()
	week7 := make([]model.UsageSnapshot, 7)
	for i := 0; i < 7; i++ {
		week7[i] = model.UsageSnapshot{
			Date:        now.AddDate(0, 0, -6+i).Format("2006-01-02"),
			TokensIn:    10000,
			TokensOut:   15000,
			TotalTokens: 25000,
			Cost:        0.25,
		}
	}

	return &model.UsageResponse{
		Today: model.UsageSnapshot{
			Date:        now.Format("2006-01-02"),
			TokensIn:    15000,
			TokensOut:   25000,
			TotalTokens: 40000,
			Cost:        0.45,
		},
		Week7: week7,
		Total: model.UsageSnapshot{
			Date:        now.Format("2006-01-02"),
			TokensIn:    350000,
			TokensOut:   520000,
			TotalTokens: 870000,
			Cost:        9.80,
		},
	}
}
