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

	sessions, err := h.client.GetSessions(ctx)
	if err != nil {
		h.logger.Warnf("获取会话列表失败: %v", err)
		// 返回模拟数据用于演示
		sessions = getMockSessions()
	}

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
	sessions, err := h.client.GetSessions(ctx)
	if err != nil {
		h.logger.Warnf("获取会话失败: %v", err)
		sessions = getMockSessions()
	}

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

	statuses, err := h.client.GetSessionStatus(ctx)
	if err != nil {
		h.logger.Warnf("获取会话状态失败: %v", err)
		statuses = getMockStatuses()
	}

	response := model.APIResponse{
		OK:   true,
		Data: statuses,
	}

	writeJSON(w, http.StatusOK, response)
}

// getMockSessions 获取模拟会话数据
// 用于 Gateway 不可用时的演示
func getMockSessions() []model.SessionSummary {
	return []model.SessionSummary{
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
	}
}

// getMockStatuses 获取模拟状态数据
func getMockStatuses() []model.SessionStatusSnapshot {
	return []model.SessionStatusSnapshot{
		{
			SessionKey: "session-001",
			Model:      "minimax-cn/MiniMax-M2.5",
			TokensIn:   15000,
			TokensOut:  25000,
			Cost:       0.35,
			UpdatedAt:  time.Now().Format(time.RFC3339),
		},
		{
			SessionKey: "session-002",
			Model:      "minimax-cn/MiniMax-M2.5",
			TokensIn:   5000,
			TokensOut:  8000,
			Cost:       0.12,
			UpdatedAt:  time.Now().Format(time.RFC3339),
		},
	}
}
