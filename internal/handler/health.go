// Package handler 提供 HTTP 请求处理功能
// 包含健康检查相关的 HTTP 处理器
package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yahao333/openclaw-go-status/internal/client"
	"github.com/yahao333/openclaw-go-status/internal/model"
)

// HealthHandler 健康检查处理器
type HealthHandler struct {
	client *client.GatewayClient // Gateway 客户端
	logger *logrus.Logger        // 日志实例
}

// NewHealthHandler 创建健康检查处理器
// 参数:
//   - client: Gateway 客户端
//   - logger: 日志实例
//
// 返回: *HealthHandler 处理器指针
func NewHealthHandler(client *client.GatewayClient, logger *logrus.Logger) *HealthHandler {
	return &HealthHandler{
		client: client,
		logger: logger,
	}
}

// Check 健康检查
// 方法: GET /health
// 返回: JSON 格式的健康状态
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	h.logger.Debugf("执行健康检查")

	checks := make(map[string]model.Check)
	status := "healthy"

	// 检查服务器
	checks["server"] = model.Check{
		Status:  "pass",
		Message: "服务器运行正常",
	}

	// 检查 Gateway 连接
	gatewayStatus := "pass"
	gatewayMessage := "Gateway 连接正常"
	if err := h.client.CheckHealth(ctx); err != nil {
		gatewayStatus = "fail"
		gatewayMessage = "Gateway 连接失败: " + err.Error()
		status = "degraded"
	}
	checks["gateway"] = model.Check{
		Status:  gatewayStatus,
		Message: gatewayMessage,
	}

	response := model.HealthResponse{
		Status:    status,
		Timestamp: time.Now().Format(time.RFC3339),
		Checks:    checks,
	}

	writeJSON(w, http.StatusOK, response)
}
