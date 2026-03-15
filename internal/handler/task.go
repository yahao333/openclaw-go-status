// Package handler 提供 HTTP 请求处理功能
// 包含任务相关的 HTTP 处理器
package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yahao333/openclaw-go-status/internal/client"
	"github.com/yahao333/openclaw-go-status/internal/model"
)

// TaskHandler 任务处理器
type TaskHandler struct {
	client *client.GatewayClient // Gateway 客户端
	logger *logrus.Logger        // 日志实例
}

// NewTaskHandler 创建任务处理器
// 参数:
//   - client: Gateway 客户端
//   - logger: 日志实例
//
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

	tasks, _ := h.client.GetTasks(ctx)

	response := model.APIResponse{
		OK:   true,
		Data: tasks,
	}

	writeJSON(w, http.StatusOK, response)
}
