// Package handler 提供 HTTP 请求处理功能
// 包含用量相关的 HTTP 处理器
package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/yahao333/openclaw-go-status/internal/client"
	"github.com/yahao333/openclaw-go-status/internal/model"
	"github.com/sirupsen/logrus"
)

// UsageHandler 用量处理器
type UsageHandler struct {
	client *client.GatewayClient // Gateway 客户端
	logger *logrus.Logger       // 日志实例
}

// NewUsageHandler 创建用量处理器
// 参数:
//   - client: Gateway 客户端
//   - logger: 日志实例
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

	usage, err := h.client.GetUsage(ctx)
	if err != nil {
		h.logger.Warnf("获取用量统计失败: %v", err)
		// 返回模拟数据用于演示
		usage = getMockUsage()
	}

	response := model.APIResponse{
		OK:   true,
		Data: usage,
	}

	writeJSON(w, http.StatusOK, response)
}

// getMockUsage 获取模拟用量数据
func getMockUsage() *model.UsageResponse {
	now := time.Now()
	week7 := make([]model.UsageSnapshot, 7)
	month30 := make([]model.UsageSnapshot, 30)

	// 生成 7 天数据
	for i := 0; i < 7; i++ {
		week7[i] = model.UsageSnapshot{
			Date:        now.AddDate(0, 0, -6+i).Format("2006-01-02"),
			TokensIn:    int64(10000 + i*1000),
			TokensOut:   int64(15000 + i*1500),
			TotalTokens: int64(25000 + i*2500),
			Cost:        0.25 + float64(i)*0.03,
		}
	}

	// 生成 30 天数据
	for i := 0; i < 30; i++ {
		month30[i] = model.UsageSnapshot{
			Date:        now.AddDate(0, 0, -29+i).Format("2006-01-02"),
			TokensIn:    int64(8000 + i*200),
			TokensOut:   int64(12000 + i*300),
			TotalTokens: int64(20000 + i*500),
			Cost:        0.20 + float64(i)*0.005,
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
		Week7:   week7,
		Month30: month30,
		Total: model.UsageSnapshot{
			Date:        now.Format("2006-01-02"),
			TokensIn:    350000,
			TokensOut:   520000,
			TotalTokens: 870000,
			Cost:        9.80,
		},
	}
}
