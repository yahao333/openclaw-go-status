// Package client 提供 Gateway 客户端功能
// 用于与 OpenClaw Gateway 通信获取状态数据
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yahao333/openclaw-go-status/internal/model"
)

// GatewayClient Gateway 客户端
// 负责与 OpenClaw Gateway HTTP API 通信
type GatewayClient struct {
	baseURL string         // Gateway 基础 URL
	client  *http.Client   // HTTP 客户端
	logger  *logrus.Logger // 日志实例
	timeout time.Duration  // 请求超时
}

// NewGatewayClient 创建 Gateway 客户端
// 参数:
//   - baseURL: Gateway 基础 URL
//   - timeout: 请求超时时间(秒)
//   - logger: 日志实例
//
// 返回: *GatewayClient 客户端指针
func NewGatewayClient(baseURL string, timeout int, logger *logrus.Logger) *GatewayClient {
	// 将 ws:// 转换为 http://
	baseURL = strings.Replace(baseURL, "ws://", "http://", 1)
	baseURL = strings.Replace(baseURL, "wss://", "https://", 1)

	return &GatewayClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
		logger:  logger,
		timeout: time.Duration(timeout) * time.Second,
	}
}

// doRequest 发送 HTTP 请求
// 参数:
//   - ctx: 上下文
//   - method: HTTP 方法
//   - path: 请求路径
//   - body: 请求体(可选)
//
// 返回: []byte 响应体, error 错误信息
func (c *GatewayClient) doRequest(ctx context.Context, method, path string, body []byte) ([]byte, error) {
	url := c.baseURL + path

	var req *http.Request
	var err error

	if len(body) > 0 {
		req, err = http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	} else {
		req, err = http.NewRequestWithContext(ctx, method, url, nil)
	}

	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	c.logger.Debugf("发送请求: %s %s", method, url)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP 错误: %d", resp.StatusCode)
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	return buf.Bytes(), nil
}

// GetSessions 获取会话列表
// 返回: []model.SessionSummary 会话列表, error 错误信息
func (c *GatewayClient) GetSessions(ctx context.Context) ([]model.SessionSummary, error) {
	data, err := c.doRequest(ctx, "GET", "/sessions", nil)
	if err != nil {
		c.logger.Warnf("获取会话列表失败: %v", err)
		return nil, err
	}

	var response model.SessionsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		c.logger.Warnf("解析会话列表失败: %v", err)
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return response.Sessions, nil
}

// GetSessionStatus 获取会话状态
// 返回: []model.SessionStatusSnapshot 状态列表, error 错误信息
func (c *GatewayClient) GetSessionStatus(ctx context.Context) ([]model.SessionStatusSnapshot, error) {
	data, err := c.doRequest(ctx, "GET", "/session-status", nil)
	if err != nil {
		c.logger.Warnf("获取会话状态失败: %v", err)
		return nil, err
	}

	var response model.SessionStatusResponse
	if err := json.Unmarshal(data, &response); err != nil {
		c.logger.Warnf("解析会话状态失败: %v", err)
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return response.Statuses, nil
}

// GetTasks 获取任务列表
// 返回: []model.ProjectTask 任务列表, error 错误信息
func (c *GatewayClient) GetTasks(ctx context.Context) ([]model.ProjectTask, error) {
	data, err := c.doRequest(ctx, "GET", "/tasks", nil)
	if err != nil {
		c.logger.Warnf("获取任务列表失败: %v", err)
		return nil, err
	}

	var response model.TasksResponse
	if err := json.Unmarshal(data, &response); err != nil {
		c.logger.Warnf("解析任务列表失败: %v", err)
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return response.Tasks, nil
}

// GetProjects 获取项目列表
// 返回: []model.ProjectRecord 项目列表, error 错误信息
func (c *GatewayClient) GetProjects(ctx context.Context) ([]model.ProjectRecord, error) {
	data, err := c.doRequest(ctx, "GET", "/projects", nil)
	if err != nil {
		c.logger.Warnf("获取项目列表失败: %v", err)
		return nil, err
	}

	var response model.ProjectsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		c.logger.Warnf("解析项目列表失败: %v", err)
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return response.Projects, nil
}

// GetCronJobs 获取 Cron 任务列表
// 返回: []model.CronJobSummary Cron 任务列表, error 错误信息
func (c *GatewayClient) GetCronJobs(ctx context.Context) ([]model.CronJobSummary, error) {
	data, err := c.doRequest(ctx, "GET", "/cron", nil)
	if err != nil {
		c.logger.Warnf("获取 Cron 任务失败: %v", err)
		return nil, err
	}

	var response model.CronResponse
	if err := json.Unmarshal(data, &response); err != nil {
		c.logger.Warnf("解析 Cron 任务失败: %v", err)
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return response.Jobs, nil
}

// GetApprovals 获取审批列表
// 返回: []model.ApprovalSummary 审批列表, error 错误信息
func (c *GatewayClient) GetApprovals(ctx context.Context) ([]model.ApprovalSummary, error) {
	data, err := c.doRequest(ctx, "GET", "/approvals", nil)
	if err != nil {
		c.logger.Warnf("获取审批列表失败: %v", err)
		return nil, err
	}

	var response model.ApprovalsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		c.logger.Warnf("解析审批列表失败: %v", err)
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return response.Approvals, nil
}

// GetExceptions 获取异常列表
// 返回: *model.ExceptionsResponse 异常响应, error 错误信息
func (c *GatewayClient) GetExceptions(ctx context.Context) (*model.ExceptionsResponse, error) {
	data, err := c.doRequest(ctx, "GET", "/exceptions", nil)
	if err != nil {
		c.logger.Warnf("获取异常列表失败: %v", err)
		return nil, err
	}

	var response model.ExceptionsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		c.logger.Warnf("解析异常列表失败: %v", err)
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &response, nil
}

// GetUsage 获取用量统计
// 返回: *model.UsageResponse 用量响应, error 错误信息
func (c *GatewayClient) GetUsage(ctx context.Context) (*model.UsageResponse, error) {
	data, err := c.doRequest(ctx, "GET", "/usage", nil)
	if err != nil {
		c.logger.Warnf("获取用量统计失败: %v", err)
		return nil, err
	}

	var response model.UsageResponse
	if err := json.Unmarshal(data, &response); err != nil {
		c.logger.Warnf("解析用量统计失败: %v", err)
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &response, nil
}

// CheckHealth 检查 Gateway 健康状态
// 返回: error 错误信息(健康时返回 nil)
func (c *GatewayClient) CheckHealth(ctx context.Context) error {
	_, err := c.doRequest(ctx, "GET", "/health", nil)
	return err
}
