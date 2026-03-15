// Package handler 提供 HTTP 请求处理功能
// 包含 Cron 任务相关的 HTTP 处理器
package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yahao333/openclaw-go-status/internal/client"
	"github.com/yahao333/openclaw-go-status/internal/model"
)

// CronHandler Cron 任务处理器
type CronHandler struct {
	client *client.GatewayClient // Gateway 客户端
	logger *logrus.Logger        // 日志实例
}

// NewCronHandler 创建 Cron 任务处理器
// 参数:
//   - client: Gateway 客户端
//   - logger: 日志实例
//
// 返回: *CronHandler 处理器指针
func NewCronHandler(client *client.GatewayClient, logger *logrus.Logger) *CronHandler {
	return &CronHandler{
		client: client,
		logger: logger,
	}
}

// List 获取 Cron 任务列表
// 方法: GET /api/cron
// 返回: JSON 格式的 Cron 任务列表
func (h *CronHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	h.logger.Infof("获取 Cron 任务列表")

	jobs, err := h.client.GetCronJobs(ctx)
	if err != nil {
		h.logger.Warnf("获取 Cron 任务列表失败: %v", err)
	}

	response := model.APIResponse{
		OK:   true,
		Data: jobs,
	}

	writeJSON(w, http.StatusOK, response)
}
