// Package model 定义数据模型
// 包含所有 API 交互中使用的数据结构
package model

import "time"

// ==================== 枚举类型 ====================

// AgentRunState Agent 运行状态
type AgentRunState string

const (
	StateIdle            AgentRunState = "idle"              // 空闲
	StateRunning         AgentRunState = "running"          // 运行中
	StateBlocked         AgentRunState = "blocked"          // 阻塞
	StateWaitingApproval AgentRunState = "waiting_approval" // 等待审批
	StateError           AgentRunState = "error"            // 错误
)

// TaskState 任务状态
type TaskState string

const (
	TaskTodo       TaskState = "todo"        // 待办
	TaskInProgress TaskState = "in_progress" // 进行中
	TaskBlocked    TaskState = "blocked"     // 阻塞
	TaskDone       TaskState = "done"        // 完成
)

// ProjectState 项目状态
type ProjectState string

const (
	ProjectPlanned ProjectState = "planned" // 计划中
	ProjectActive  ProjectState = "active"  // 进行中
	ProjectBlocked ProjectState = "blocked" // 阻塞
	ProjectDone    ProjectState = "done"    // 完成
)

// BudgetStatus 预算状态
type BudgetStatus string

const (
	BudgetOk   BudgetStatus = "ok"   // 正常
	BudgetWarn BudgetStatus = "warn" // 警告
	BudgetOver BudgetStatus = "over" // 超预算
)

// ApprovalState 审批状态
type ApprovalState string

const (
	ApprovalPending ApprovalState = "pending"  // 待审批
	ApprovalApproved ApprovalState = "approved" // 已批准
	ApprovalDenied  ApprovalState = "denied"   // 已拒绝
	ApprovalUnknown ApprovalState = "unknown" // 未知
)

// AlertLevel 告警级别
type AlertLevel string

const (
	AlertInfo          AlertLevel = "info"           // 信息
	AlertWarn          AlertLevel = "warn"           // 警告
	AlertActionRequired AlertLevel = "action-required" // 需要操作
)

// ==================== 会话相关 ====================

// SessionSummary 会话摘要
type SessionSummary struct {
	SessionKey  string       `json:"sessionKey"`  // 会话唯一标识
	Label       string       `json:"label,omitempty"` // 会话标签
	AgentID     string       `json:"agentId,omitempty"` // Agent ID
	State       AgentRunState `json:"state"`        // 运行状态
	LastMessageAt string     `json:"lastMessageAt,omitempty"` // 最后消息时间
}

// SessionStatusSnapshot 会话状态快照
type SessionStatusSnapshot struct {
	SessionKey  string  `json:"sessionKey"` // 会话唯一标识
	Model       string  `json:"model,omitempty"` // 使用模型
	TokensIn    int64   `json:"tokensIn,omitempty"`  // 输入 Token 数
	TokensOut   int64   `json:"tokensOut,omitempty"` // 输出 Token 数
	Cost        float64 `json:"cost,omitempty"`     // 费用
	UpdatedAt   string  `json:"updatedAt"`          // 更新时间
}

// SessionsResponse 会话列表响应
type SessionsResponse struct {
	Sessions []SessionSummary `json:"sessions"` // 会话列表
}

// SessionStatusResponse 会话状态响应
type SessionStatusResponse struct {
	Statuses []SessionStatusSnapshot `json:"statuses"` // 状态列表
}

// ==================== 任务相关 ====================

// TaskArtifact 任务产物
type TaskArtifact struct {
	ArtifactID string `json:"artifactId"` // 产物 ID
	Type       string `json:"type"`       // 类型: code, doc, link, other
	Label      string `json:"label"`      // 标签
	Location   string `json:"location"`  // 位置
}

// BudgetThresholds 预算阈值
type BudgetThresholds struct {
	TokensIn   *int64   `json:"tokensIn,omitempty"`   // 输入 Token 上限
	TokensOut  *int64   `json:"tokensOut,omitempty"`  // 输出 Token 上限
	TotalTokens *int64  `json:"totalTokens,omitempty"` // 总 Token 上限
	Cost       *float64 `json:"cost,omitempty"`       // 费用上限
	WarnRatio  *float64 `json:"warnRatio,omitempty"` // 警告比例
}

// ProjectTask 项目任务
type ProjectTask struct {
	ProjectID   string         `json:"projectId"`   // 项目 ID
	TaskID      string         `json:"taskId"`      // 任务 ID
	Title       string         `json:"title"`       // 标题
	Status      TaskState      `json:"status"`      // 状态
	Owner       string         `json:"owner"`       // 负责人
	DueAt       *time.Time     `json:"dueAt,omitempty"` // 截止时间
	Artifacts   []TaskArtifact `json:"artifacts"`   // 产物列表
	SessionKeys []string       `json:"sessionKeys"` // 关联会话
	Budget      BudgetThresholds `json:"budget"`    // 预算
	UpdatedAt   string         `json:"updatedAt"`   // 更新时间
}

// TasksResponse 任务列表响应
type TasksResponse struct {
	Tasks []ProjectTask `json:"tasks"` // 任务列表
}

// ==================== 项目相关 ====================

// ProjectRecord 项目记录
type ProjectRecord struct {
	ProjectID string       `json:"projectId"` // 项目 ID
	Title     string       `json:"title"`     // 标题
	Status    ProjectState `json:"status"`    // 状态
	Owner     string       `json:"owner"`     // 负责人
	Budget    BudgetThresholds `json:"budget"` // 预算
	UpdatedAt string       `json:"updatedAt"` // 更新时间
}

// ProjectsResponse 项目列表响应
type ProjectsResponse struct {
	Projects []ProjectRecord `json:"projects"` // 项目列表
}

// ==================== 用量相关 ====================

// UsageSnapshot 用量快照
type UsageSnapshot struct {
	Date        string  `json:"date"`         // 日期
	TokensIn    int64   `json:"tokensIn"`    // 输入 Token
	TokensOut   int64   `json:"tokensOut"`    // 输出 Token
	TotalTokens int64   `json:"totalTokens"` // 总 Token
	Cost        float64 `json:"cost"`        // 费用
}

// UsageResponse 用量响应
type UsageResponse struct {
	Today    UsageSnapshot  `json:"today"`     // 今日用量
	Week7   []UsageSnapshot `json:"week7"`     // 7天用量
	Month30 []UsageSnapshot `json:"month30"`   // 30天用量
	Total   UsageSnapshot   `json:"total"`     // 汇总
}

// ==================== Cron 相关 ====================

// CronJobSummary Cron 任务摘要
type CronJobSummary struct {
	JobID     string  `json:"jobId"`      // 任务 ID
	Name      string  `json:"name,omitempty"` // 名称
	Enabled   bool    `json:"enabled"`    // 是否启用
	NextRunAt string  `json:"nextRunAt,omitempty"` // 下次运行时间
	Health    string  `json:"health"`     // 健康状态: scheduled, due, late, unknown, disabled
}

// CronResponse Cron 响应
type CronResponse struct {
	Jobs []CronJobSummary `json:"jobs"` // 任务列表
}

// ==================== 审批相关 ====================

// ApprovalSummary 审批摘要
type ApprovalSummary struct {
	ApprovalID string        `json:"approvalId"` // 审批 ID
	SessionKey string        `json:"sessionKey,omitempty"` // 会话 ID
	AgentID    string        `json:"agentId,omitempty"` // Agent ID
	Status     ApprovalState `json:"status"`     // 状态
	Decision   string        `json:"decision,omitempty"` // 决策
	Command    string        `json:"command,omitempty"` // 命令
	Reason     string        `json:"reason,omitempty"` // 原因
	RequestedAt string       `json:"requestedAt,omitempty"` // 请求时间
	UpdatedAt  string        `json:"updatedAt,omitempty"` // 更新时间
}

// ApprovalsResponse 审批列表响应
type ApprovalsResponse struct {
	Approvals []ApprovalSummary `json:"approvals"` // 审批列表
}

// ==================== 异常相关 ====================

// ExceptionItem 异常项
type ExceptionItem struct {
	Level     AlertLevel `json:"level"`      // 级别
	Code      string     `json:"code"`       // 错误码
	Source    string     `json:"source"`     // 来源: system, session, approval, budget, task
	SourceID  string     `json:"sourceId"`   // 源 ID
	Message   string     `json:"message"`    // 消息
	Route     string     `json:"route"`      // 路由: timeline, operator-watch, action-queue
	OccurredAt string    `json:"occurredAt,omitempty"` // 发生时间
}

// ExceptionsResponse 异常列表响应
type ExceptionsResponse struct {
	GeneratedAt string          `json:"generatedAt"` // 生成时间
	Items      []ExceptionItem `json:"items"`      // 异常列表
	Counts     ExceptionCounts `json:"counts"`     // 统计
}

// ExceptionCounts 异常统计
type ExceptionCounts struct {
	Info          int `json:"info"`           // 信息级
	Warn          int `json:"warn"`           // 警告级
	ActionRequired int `json:"actionRequired"` // 需要操作
}

// ==================== 健康检查 ====================

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status    string            `json:"status"`     // 状态: healthy, degraded, unhealthy
	Timestamp string            `json:"timestamp"` // 时间戳
	Checks    map[string]Check `json:"checks"`    // 检查项
}

// Check 检查项
type Check struct {
	Status  string `json:"status"`  // 状态: pass, fail
	Message string `json:"message"` // 消息
}

// ==================== 通用响应 ====================

// ErrorResponse 错误响应
type ErrorResponse struct {
	Error     string `json:"error"`      // 错误信息
	Code      string `json:"code,omitempty"` // 错误码
	RequestID string `json:"requestId,omitempty"` // 请求 ID
}

// APIResponse 通用 API 响应
type APIResponse struct {
	OK      bool        `json:"ok"`       // 是否成功
	Data    interface{} `json:"data,omitempty"` // 数据
	Error   string      `json:"error,omitempty"` // 错误信息
	Message string      `json:"message,omitempty"` // 消息
}
