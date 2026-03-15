// Package handler 提供 HTTP 请求处理功能
// 包含会话、任务、项目等相关的 HTTP 处理器
package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yahao333/openclaw-go-status/internal/client"
	"github.com/yahao333/openclaw-go-status/internal/model"
)

// SessionHandler 会话处理器
type SessionHandler struct {
	client *client.GatewayClient // Gateway 客户端
	logger *logrus.Logger        // 日志实例
}

// NewSessionHandler 创建会话处理器
// 参数:
//   - client: Gateway 客户端
//   - logger: 日志实例
//
// 返回: *SessionHandler 处理器指针
func NewSessionHandler(client *client.GatewayClient, logger *logrus.Logger) *SessionHandler {
	return &SessionHandler{
		client: client,
		logger: logger,
	}
}

// List 获取会话列表
// 方法: GET /api/sessions
// 返回: JSON 格式的会话列表
func (h *SessionHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	h.logger.Infof("获取会话列表")

	sessions, _ := h.client.GetSessions(ctx)

	response := model.APIResponse{
		OK:   true,
		Data: sessions,
	}

	writeJSON(w, http.StatusOK, response)
}

// Get 获取指定会话详情
// 方法: GET /api/sessions/:id
// 参数:
//   - id: 会话 ID
//
// 返回: JSON 格式的会话详情
func (h *SessionHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	sessionID := extractID(r)
	h.logger.Infof("获取会话详情: %s", sessionID)

	// 获取会话列表
	sessions, _ := h.client.GetSessions(ctx)

	// 查找指定会话
	var session *model.SessionSummary
	for i := range sessions {
		if sessions[i].SessionKey == sessionID {
			session = &sessions[i]
			break
		}
	}

	if session == nil {
		h.logger.Warnf("会话不存在: %s", sessionID)
		writeError(w, http.StatusNotFound, "会话不存在")
		return
	}

	response := model.APIResponse{
		OK:   true,
		Data: session,
	}

	writeJSON(w, http.StatusOK, response)
}

// Status 获取会话状态
// 方法: GET /api/status
// 返回: JSON 格式的会话状态列表
func (h *SessionHandler) Status(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	h.logger.Infof("获取会话状态")

	statuses, _ := h.client.GetSessionStatus(ctx)

	response := model.APIResponse{
		OK:   true,
		Data: statuses,
	}

	writeJSON(w, http.StatusOK, response)
}

// DashboardStats 获取 Dashboard 统计数据
// 方法: GET /api/dashboard
// 返回: JSON 格式的聚合统计数据
func (h *SessionHandler) DashboardStats(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	h.logger.Infof("获取 Dashboard 统计")

	// 获取会话列表和状态
	sessions, _ := h.client.GetSessions(ctx)
	statuses, _ := h.client.GetSessionStatus(ctx)

	// 构建 sessionKey -> status 映射
	statusMap := make(map[string]model.SessionStatusSnapshot)
	for _, s := range statuses {
		statusMap[s.SessionKey] = s
	}

	// 合并会话和状态数据
	mergedSessions := make([]model.SessionWithStats, len(sessions))
	runningCount := 0
	for i, s := range sessions {
		stats := statusMap[s.SessionKey]
		mergedSessions[i] = model.SessionWithStats{
			SessionSummary: s,
			TokensIn:       stats.TokensIn,
			TokensOut:      stats.TokensOut,
			Cost:           stats.Cost,
		}
		if s.State == model.StateRunning {
			runningCount++
		}
	}

	// 计算任务和项目数量（目前从会话推断）
	taskCount := 0
	projectCount := 0

	stats := model.DashboardStatsResponse{
		Sessions: len(mergedSessions),
		Running:  runningCount,
		Tasks:    taskCount,
		Projects: projectCount,
		Data:     mergedSessions,
	}

	response := model.APIResponse{
		OK:   true,
		Data: stats,
	}

	writeJSON(w, http.StatusOK, response)
}
