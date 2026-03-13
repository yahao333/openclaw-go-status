// Package handler 提供 HTTP 请求处理功能
// 包含模板渲染相关的处理函数，从 OpenClaw Gateway 获取真实数据
package handler

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/yahao333/openclaw-go-status/internal/client"
	"github.com/yahao333/openclaw-go-status/internal/model"
	"github.com/sirupsen/logrus"
)

// TemplateHandler 模板处理器
// 从 OpenClaw Gateway 获取真实数据并渲染页面
type TemplateHandler struct {
	templates map[string]*template.Template // 每个页面独立的模板集合
	client    *client.GatewayClient         // Gateway 客户端
	logger    *logrus.Logger                // 日志实例
}

// NewTemplateHandler 创建模板处理器
// 参数:
//   - templateDir: 模板目录路径
//   - gatewayClient: Gateway 客户端
//   - logger: 日志实例
//
// 返回: *TemplateHandler 处理器指针
func NewTemplateHandler(templateDir string, gatewayClient *client.GatewayClient, logger *logrus.Logger) (*TemplateHandler, error) {
	// 检查模板目录是否存在
	if _, err := os.Stat(templateDir); err != nil {
		return nil, fmt.Errorf("模板目录不存在: %s, 错误: %v", templateDir, err)
	}

	funcMap := template.FuncMap{
		"formatNumber": formatNumber,
		"formatTime":   formatTime,
	}

	basePath := filepath.Join(templateDir, "base.html")
	
	// 检查 base.html 是否存在
	if _, err := os.Stat(basePath); err != nil {
		return nil, fmt.Errorf("base.html 不存在: %s, 错误: %v", basePath, err)
	}

	// 获取所有 HTML 文件
	files, err := filepath.Glob(filepath.Join(templateDir, "*.html"))
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("模板目录 %s 中未找到 HTML 文件", templateDir)
	}

	templates := make(map[string]*template.Template)
	for _, filePath := range files {
		baseName := filepath.Base(filePath)
		if baseName == "base.html" {
			continue
		}

		// 使用 template.Must 确保模板解析成功
		tmpl := template.Must(template.New(baseName).Funcs(funcMap).ParseFiles(basePath, filePath))
		templates[baseName] = tmpl
		logger.Infof("已加载模板: %s", baseName)
	}

	if len(templates) == 0 {
		return nil, fmt.Errorf("模板目录 %s 中未找到页面模板", templateDir)
	}

	logger.Infof("共加载 %d 个页面模板", len(templates))

	return &TemplateHandler{
		templates: templates,
		client:    gatewayClient,
		logger:    logger,
	}, nil
}

// Home 首页/仪表盘
// 方法: GET /
// 返回: HTML 页面
func (h *TemplateHandler) Home(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	// 获取健康状态
	healthStatus, healthMessage := h.getHealthStatus(ctx)

	// 获取会话列表
	sessions, _ := h.client.GetSessions(ctx)

	// 获取会话状态
	statuses, _ := h.client.GetSessionStatus(ctx)

	// 合并会话和状态数据
	sessionMap := make(map[string]model.SessionStatusSnapshot)
	for _, s := range statuses {
		sessionMap[s.SessionKey] = s
	}

	sessionsWithStatus := make([]SessionWithStatus, 0)
	runningCount := 0
	for _, s := range sessions {
		status := sessionMap[s.SessionKey]
		sessionsWithStatus = append(sessionsWithStatus, SessionWithStatus{
			SessionSummary: s,
			TokensIn:      status.TokensIn,
			TokensOut:     status.TokensOut,
			Cost:          status.Cost,
		})
		if s.State == model.StateRunning {
			runningCount++
		}
	}

	// 获取任务列表
	tasks, _ := h.client.GetTasks(ctx)
	todoCount, inProgressCount, blockedCount, doneCount := countTasks(tasks)

	// 获取项目列表
	projects, _ := h.client.GetProjects(ctx)

	// 获取用量
	usage, _ := h.client.GetUsage(ctx)

	data := PageData{
		ActivePage:    "dashboard",
		LastUpdate:    time.Now().Format("2006-01-02 15:04:05"),
		HealthStatus:  healthStatus,
		HealthMessage: healthMessage,
		Stats: Stats{
			Sessions: len(sessions),
			Running:  runningCount,
			Tasks:    len(tasks),
			Projects: len(projects),
		},
		Sessions:       sessionsWithStatus,
		RecentSessions: sessions,
		Tasks:          tasks,
		RecentTasks:    tasks,
		Projects:       projects,
		TaskStats: TaskStats{
			Todo:       todoCount,
			InProgress: inProgressCount,
			Blocked:    blockedCount,
			Done:       doneCount,
		},
		Usage: usage,
	}

	if usage != nil {
		data.TodayUsage = usage.Today
	}

	h.render(w, "index.html", data)
}

// Sessions 会话列表页
// 方法: GET /sessions
// 返回: HTML 页面
func (h *TemplateHandler) Sessions(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	healthStatus, healthMessage := h.getHealthStatus(ctx)
	sessions, _ := h.client.GetSessions(ctx)
	statuses, _ := h.client.GetSessionStatus(ctx)

	sessionMap := make(map[string]model.SessionStatusSnapshot)
	for _, s := range statuses {
		sessionMap[s.SessionKey] = s
	}

	sessionsWithStatus := make([]SessionWithStatus, 0)
	for _, s := range sessions {
		status := sessionMap[s.SessionKey]
		sessionsWithStatus = append(sessionsWithStatus, SessionWithStatus{
			SessionSummary: s,
			TokensIn:      status.TokensIn,
			TokensOut:     status.TokensOut,
			Cost:          status.Cost,
		})
	}

	data := PageData{
		ActivePage:    "sessions",
		LastUpdate:    time.Now().Format("2006-01-02 15:04:05"),
		HealthStatus:  healthStatus,
		HealthMessage: healthMessage,
		Sessions:      sessionsWithStatus,
	}

	h.render(w, "sessions.html", data)
}

// Tasks 任务列表页
// 方法: GET /tasks
// 返回: HTML 页面
func (h *TemplateHandler) Tasks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	healthStatus, healthMessage := h.getHealthStatus(ctx)
	tasks, _ := h.client.GetTasks(ctx)

	todoCount, inProgressCount, blockedCount, doneCount := countTasks(tasks)

	data := PageData{
		ActivePage:    "tasks",
		LastUpdate:    time.Now().Format("2006-01-02 15:04:05"),
		HealthStatus:  healthStatus,
		HealthMessage: healthMessage,
		Tasks:         tasks,
		TaskStats: TaskStats{
			Todo:       todoCount,
			InProgress: inProgressCount,
			Blocked:    blockedCount,
			Done:       doneCount,
		},
	}

	h.render(w, "tasks.html", data)
}

// Projects 项目列表页
// 方法: GET /projects
// 返回: HTML 页面
func (h *TemplateHandler) Projects(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	healthStatus, healthMessage := h.getHealthStatus(ctx)
	projects, _ := h.client.GetProjects(ctx)

	data := PageData{
		ActivePage:    "projects",
		LastUpdate:    time.Now().Format("2006-01-02 15:04:05"),
		HealthStatus:  healthStatus,
		HealthMessage: healthMessage,
		Projects:      projects,
	}

	h.render(w, "projects.html", data)
}

// Usage 用量统计页
// 方法: GET /usage
// 返回: HTML 页面
func (h *TemplateHandler) Usage(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	healthStatus, healthMessage := h.getHealthStatus(ctx)
	usage, _ := h.client.GetUsage(ctx)

	data := PageData{
		ActivePage:    "usage",
		LastUpdate:    time.Now().Format("2006-01-02 15:04:05"),
		HealthStatus:  healthStatus,
		HealthMessage: healthMessage,
		Usage:         usage,
	}

	h.render(w, "usage.html", data)
}

// getHealthStatus 获取健康状态
func (h *TemplateHandler) getHealthStatus(ctx context.Context) (string, string) {
	if err := h.client.CheckHealth(ctx); err != nil {
		h.logger.Warnf("Gateway 健康检查失败: %v", err)
		return "unhealthy", "Gateway 连接失败"
	}
	return "healthy", "Gateway 连接正常"
}

// countTasks 统计任务数量
func countTasks(tasks []model.ProjectTask) (todo, inProgress, blocked, done int) {
	for _, task := range tasks {
		switch task.Status {
		case model.TaskTodo:
			todo++
		case model.TaskInProgress:
			inProgress++
		case model.TaskBlocked:
			blocked++
		case model.TaskDone:
			done++
		}
	}
	return
}

// render 渲染模板
func (h *TemplateHandler) render(w http.ResponseWriter, name string, data interface{}) {
	tmpl, ok := h.templates[name]
	if !ok {
		h.logger.Errorf("模板不存在: %s", name)
		http.Error(w, fmt.Sprintf("模板不存在: %s", name), http.StatusInternalServerError)
		return
	}
	
	if tmpl == nil {
		h.logger.Errorf("模板为空: %s", name)
		http.Error(w, fmt.Sprintf("模板为空: %s", name), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		h.logger.Errorf("渲染模板 %s 失败: %v", name, err)
		http.Error(w, "模板渲染失败: "+err.Error(), http.StatusInternalServerError)
	}
}

// ==================== 辅助函数 ====================

// formatNumber 格式化数字（添加千分位分隔符）
func formatNumber(n int64) string {
	if n == 0 {
		return "0"
	}

	str := ""
	count := 0
	for n > 0 {
		if count > 0 && count%3 == 0 {
			str = "," + str
		}
		str = string(rune('0'+n%10)) + str
		n /= 10
		count++
	}
	return str
}

// formatTime 格式化时间
func formatTime(t *time.Time) string {
	if t == nil {
		return "-"
	}
	return t.Format("2006-01-02 15:04")
}

// ==================== 页面数据结构 ====================

// PageData 页面通用数据
type PageData struct {
	ActivePage     string                // 当前页面标识
	LastUpdate     string                // 最后更新时间
	HealthStatus   string                // 健康状态
	HealthMessage  string                // 健康消息
	Stats          Stats                 // 统计信息
	Sessions       []SessionWithStatus   // 会话列表(带状态)
	RecentSessions []model.SessionSummary // 最近会话
	Tasks          []model.ProjectTask   // 任务列表
	RecentTasks    []model.ProjectTask   // 最近任务
	Projects       []model.ProjectRecord // 项目列表
	TaskStats      TaskStats             // 任务统计
	Usage          *model.UsageResponse // 用量数据
	TodayUsage     model.UsageSnapshot   // 今日用量
}

// Stats 统计数据
type Stats struct {
	Sessions int // 会话数
	Running  int // 运行中
	Tasks    int // 任务数
	Projects int // 项目数
}

// TaskStats 任务统计
type TaskStats struct {
	Todo       int // 待办
	InProgress int // 进行中
	Blocked    int // 阻塞
	Done       int // 已完成
}

// SessionWithStatus 会话状态信息
type SessionWithStatus struct {
	model.SessionSummary
	TokensIn  int64   `json:"tokensIn"`  // 输入 Token
	TokensOut int64   `json:"tokensOut"` // 输出 Token
	Cost      float64 `json:"cost"`      // 费用
}
