// Package handler 提供 HTTP 请求处理功能
// 包含用量相关的 HTTP 处理器
package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yahao333/openclaw-go-status/internal/client"
	"github.com/yahao333/openclaw-go-status/internal/model"
)

// UsageHandler 用量处理器
type UsageHandler struct {
	client *client.GatewayClient // Gateway 客户端
	logger *logrus.Logger        // 日志实例
}

// NewUsageHandler 创建用量处理器
// 参数:
//   - client: Gateway 客户端
//   - logger: 日志实例
//
// 返回: *UsageHandler 处理器指针
func NewUsageHandler(client *client.GatewayClient, logger *logrus.Logger) *UsageHandler {
	return &UsageHandler{
		client: client,
		logger: logger,
	}
}

// Get 获取用量统计
// 方法: GET /api/usage
// 返回: JSON 格式的用量统计
func (h *UsageHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	h.logger.Infof("获取用量统计")

	usage, _ := h.client.GetUsage(ctx)

	response := model.APIResponse{
		OK:   true,
		Data: usage,
	}

	writeJSON(w, http.StatusOK, response)
}
