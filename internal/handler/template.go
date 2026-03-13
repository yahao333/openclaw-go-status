// Package handler 提供 HTTP 请求处理功能
// 包含模板渲染相关的处理函数
package handler

import (
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/yahao333/openclaw-go-status/internal/model"
)

// TemplateHandler 模板处理器
type TemplateHandler struct {
	templates *template.Template // HTML 模板
}

// NewTemplateHandler 创建模板处理器
// 参数:
//   - templateDir: 模板目录路径
//
// 返回: *TemplateHandler 处理器指针
func NewTemplateHandler(templateDir string) (*TemplateHandler, error) {
	funcMap := template.FuncMap{
		"formatNumber": formatNumber,
		"formatTime":   formatTime,
	}

	templates, err := template.New("").Funcs(funcMap).ParseGlob(filepath.Join(templateDir, "*.html"))
	if err != nil {
		return nil, err
	}

	return &TemplateHandler{
		templates: templates,
	}, nil
}

// Home 首页/仪表盘
// 方法: GET /
// 返回: HTML 页面
func (h *TemplateHandler) Home(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		ActivePage:    "dashboard",
		LastUpdate:    time.Now().Format("2006-01-02 15:04:05"),
		HealthStatus:  "healthy",
		HealthMessage: "Gateway 连接正常",
		Stats: Stats{
			Sessions: 2,
			Running:  1,
			Tasks:    3,
			Projects: 3,
		},
		RecentSessions: []model.SessionSummary{
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
		},
		RecentTasks: []model.ProjectTask{
			{
				ProjectID: "project-001",
				TaskID:    "task-001",
				Title:     "完成用户认证模块",
				Status:    model.TaskInProgress,
				Owner:     "zhangsan",
				UpdatedAt: time.Now().Format(time.RFC3339),
			},
			{
				ProjectID: "project-001",
				TaskID:    "task-002",
				Title:     "编写 API 文档",
				Status:    model.TaskTodo,
				Owner:     "lisi",
				UpdatedAt: time.Now().Format(time.RFC3339),
			},
		},
		TodayUsage: model.UsageSnapshot{
			Date:        time.Now().Format("2006-01-02"),
			TokensIn:    15000,
			TokensOut:   25000,
			TotalTokens: 40000,
			Cost:        0.45,
		},
	}

	h.render(w, "index.html", data)
}

// Sessions 会话列表页
// 方法: GET /sessions
// 返回: HTML 页面
func (h *TemplateHandler) Sessions(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		ActivePage:    "sessions",
		LastUpdate:    time.Now().Format("2006-01-02 15:04:05"),
		HealthStatus:  "healthy",
		HealthMessage: "Gateway 连接正常",
		Sessions: []SessionWithStatus{
			{
				SessionSummary: model.SessionSummary{
					SessionKey:    "session-001",
					Label:         "主会话",
					AgentID:       "agent-001",
					State:         model.StateRunning,
					LastMessageAt: time.Now().Format(time.RFC3339),
				},
				TokensIn:  15000,
				TokensOut: 25000,
				Cost:      0.35,
			},
			{
				SessionSummary: model.SessionSummary{
					SessionKey:    "session-002",
					Label:         "辅助会话",
					AgentID:       "agent-002",
					State:         model.StateIdle,
					LastMessageAt: time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
				},
				TokensIn:  5000,
				TokensOut: 8000,
				Cost:      0.12,
			},
		},
	}

	h.render(w, "sessions.html", data)
}

// Tasks 任务列表页
// 方法: GET /tasks
// 返回: HTML 页面
func (h *TemplateHandler) Tasks(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		ActivePage:    "tasks",
		LastUpdate:    time.Now().Format("2006-01-02 15:04:05"),
		HealthStatus:  "healthy",
		HealthMessage: "Gateway 连接正常",
		Tasks: []model.ProjectTask{
			{
				ProjectID:   "project-001",
				TaskID:      "task-001",
				Title:       "完成用户认证模块",
				Status:      model.TaskInProgress,
				Owner:       "zhangsan",
				SessionKeys: []string{"session-001"},
				UpdatedAt:   time.Now().Format(time.RFC3339),
			},
			{
				ProjectID:   "project-001",
				TaskID:      "task-002",
				Title:       "编写 API 文档",
				Status:      model.TaskTodo,
				Owner:       "lisi",
				SessionKeys: []string{},
				UpdatedAt:   time.Now().Format(time.RFC3339),
			},
			{
				ProjectID:   "project-001",
				TaskID:      "task-003",
				Title:       "修复登录 Bug",
				Status:      model.TaskDone,
				Owner:       "zhangsan",
				SessionKeys: []string{"session-002"},
				UpdatedAt:   time.Now().Format(time.RFC3339),
			},
		},
		TaskStats: TaskStats{
			Todo:       1,
			InProgress: 1,
			Blocked:    0,
			Done:       1,
		},
	}

	h.render(w, "tasks.html", data)
}

// Projects 项目列表页
// 方法: GET /projects
// 返回: HTML 页面
func (h *TemplateHandler) Projects(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		ActivePage:    "projects",
		LastUpdate:    time.Now().Format("2006-01-02 15:04:05"),
		HealthStatus:  "healthy",
		HealthMessage: "Gateway 连接正常",
		Projects: []model.ProjectRecord{
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
			{
				ProjectID: "project-003",
				Title:     "系统优化",
				Status:    model.ProjectDone,
				Owner:     "wangwu",
				UpdatedAt: time.Now().Format(time.RFC3339),
			},
		},
	}

	h.render(w, "projects.html", data)
}

// Usage 用量统计页
// 方法: GET /usage
// 返回: HTML 页面
func (h *TemplateHandler) Usage(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	week7 := make([]model.UsageSnapshot, 7)

	for i := 0; i < 7; i++ {
		week7[i] = model.UsageSnapshot{
			Date:        now.AddDate(0, 0, -6+i).Format("2006-01-02"),
			TokensIn:    int64(10000 + i*1000),
			TokensOut:   int64(15000 + i*1500),
			TotalTokens: int64(25000 + i*2500),
			Cost:        0.25 + float64(i)*0.03,
		}
	}

	data := PageData{
		ActivePage:    "usage",
		LastUpdate:    time.Now().Format("2006-01-02 15:04:05"),
		HealthStatus:  "healthy",
		HealthMessage: "Gateway 连接正常",
		Usage: &model.UsageResponse{
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
		},
	}

	h.render(w, "usage.html", data)
}

// render 渲染模板
// 参数:
//   - w: HTTP 响应写入器
//   - name: 模板名称
//   - data: 页面数据
func (h *TemplateHandler) render(w http.ResponseWriter, name string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, name, data); err != nil {
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
	ActivePage     string                 // 当前页面标识
	LastUpdate     string                 // 最后更新时间
	HealthStatus   string                 // 健康状态
	HealthMessage  string                 // 健康消息
	Stats          Stats                  // 统计信息
	Sessions       []SessionWithStatus    // 会话列表(带状态)
	RecentSessions []model.SessionSummary // 最近会话
	Tasks          []model.ProjectTask    // 任务列表
	RecentTasks    []model.ProjectTask    // 最近任务
	Projects       []model.ProjectRecord  // 项目列表
	TaskStats      TaskStats              // 任务统计
	Usage          *model.UsageResponse   // 用量数据
	TodayUsage     model.UsageSnapshot    // 今日用量
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
