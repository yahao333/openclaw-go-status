// Package handler 提供 HTTP 请求处理功能
// 包含任务相关的 HTTP 处理器
package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/yahao333/openclaw-go-status/internal/client"
	"github.com/yahao333/openclaw-go-status/internal/model"
	"github.com/sirupsen/logrus"
)

// TaskHandler 任务处理器
type TaskHandler struct {
	client *client.GatewayClient // Gateway 客户端
	logger *logrus.Logger       // 日志实例
}

// NewTaskHandler 创建任务处理器
// 参数:
//   - client: Gateway 客户端
//   - logger: 日志实例
// 返回: *TaskHandler 处理器指针
func NewTaskHandler(client *client.GatewayClient, logger *logrus.Logger) *TaskHandler {
	return &TaskHandler{
		client: client,
		logger: logger,
	}
}

// List 获取任务列表
// 方法: GET /api/tasks
// 返回: JSON 格式的任务列表
func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	h.logger.Infof("获取任务列表")

	tasks, err := h.client.GetTasks(ctx)
	if err != nil {
		h.logger.Warnf("获取任务列表失败: %v", err)
		// 返回模拟数据用于演示
		tasks = getMockTasks()
	}

	response := model.APIResponse{
		OK:   true,
		Data: tasks,
	}

	writeJSON(w, http.StatusOK, response)
}

// getMockTasks 获取模拟任务数据
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
	}
}
