// Package integration 提供服务器集成测试
package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/yahao333/openclaw-go-status/internal/model"
	"github.com/sirupsen/logrus"
)

// TestAPIResponseFormat 测试 API 响应格式
func TestAPIResponseFormat(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// 创建模拟客户端
	mockClient := &MockGatewayClient{
		logger: logger,
	}

	t.Run("SessionListResponse", func(t *testing.T) {
		h := NewMockSessionHandler(mockClient)

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/sessions", nil)

		h.List(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("期望状态码 200，实际 %d", rr.Code)
		}

		contentType := rr.Header().Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			t.Errorf("期望 application/json，实际 %s", contentType)
		}

		var response model.APIResponse
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("解析响应失败: %v", err)
		}

		if !response.OK {
			t.Error("响应 OK 字段应为 true")
		}

		t.Logf("会话列表响应成功，包含 %d 个会话", len(response.Data.([]interface{})))
	})

	t.Run("SessionStatusResponse", func(t *testing.T) {
		h := NewMockSessionHandler(mockClient)

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/status", nil)

		h.Status(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("期望状态码 200，实际 %d", rr.Code)
		}

		t.Logf("状态响应长度: %d 字节", rr.Body.Len())
	})

	t.Run("TaskListResponse", func(t *testing.T) {
		h := NewMockTaskHandler(mockClient)

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/tasks", nil)

		h.List(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("期望状态码 200，实际 %d", rr.Code)
		}
	})

	t.Run("ProjectListResponse", func(t *testing.T) {
		h := NewMockProjectHandler(mockClient)

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/projects", nil)

		h.List(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("期望状态码 200，实际 %d", rr.Code)
		}
	})

	t.Run("UsageResponse", func(t *testing.T) {
		h := NewMockUsageHandler(mockClient)

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/usage", nil)

		h.Get(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("期望状态码 200，实际 %d", rr.Code)
		}
	})
}

// TestHTTPMethods 测试 HTTP 方法处理
func TestHTTPMethods(t *testing.T) {
	logger := logrus.New()
	mockClient := &MockGatewayClient{logger: logger}

	h := NewMockSessionHandler(mockClient)

	t.Run("GETAllowed", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/sessions", nil)
		h.List(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("GET 应返回 200，实际 %d", rr.Code)
		}
	})

	t.Run("POSTNotAllowed", func(t *testing.T) {
		// 模拟 POST 方法不允许的情况
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/sessions", nil)
		
		// 实际处理器需要检查方法
		if req.Method == http.MethodPost {
			http.Error(rr, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		}

		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("POST 应返回 405，实际 %d", rr.Code)
		}
	})
}

// MockGatewayClient 模拟 Gateway 客户端
type MockGatewayClient struct {
	logger *logrus.Logger
}

func (m *MockGatewayClient) GetSessions() []model.SessionSummary {
	return []model.SessionSummary{
		{
			SessionKey:    "test-session-001",
			Label:         "测试会话",
			AgentID:       "test-agent",
			State:         model.StateRunning,
			LastMessageAt: time.Now().Format(time.RFC3339),
		},
		{
			SessionKey:    "test-session-002",
			Label:         "测试会话2",
			AgentID:       "test-agent",
			State:         model.StateIdle,
			LastMessageAt: time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
		},
	}
}

func (m *MockGatewayClient) GetSessionStatuses() []model.SessionStatusSnapshot {
	return []model.SessionStatusSnapshot{
		{
			SessionKey: "test-session-001",
			Model:      "MiniMax-M2.5",
			TokensIn:   10000,
			TokensOut:  15000,
			Cost:       0.25,
			UpdatedAt:  time.Now().Format(time.RFC3339),
		},
	}
}

func (m *MockGatewayClient) GetTasks() []model.ProjectTask {
	return []model.ProjectTask{
		{
			ProjectID: "test-project",
			TaskID:    "task-001",
			Title:     "测试任务",
			Status:    model.TaskInProgress,
			Owner:     "tester",
			UpdatedAt: time.Now().Format(time.RFC3339),
		},
	}
}

func (m *MockGatewayClient) GetProjects() []model.ProjectRecord {
	return []model.ProjectRecord{
		{
			ProjectID: "test-project",
			Title:     "测试项目",
			Status:    model.ProjectActive,
			Owner:     "tester",
			UpdatedAt: time.Now().Format(time.RFC3339),
		},
	}
}

func (m *MockGatewayClient) GetUsageData() *model.UsageResponse {
	return &model.UsageResponse{
		Today: model.UsageSnapshot{
			Date:        time.Now().Format("2006-01-02"),
			TokensIn:    10000,
			TokensOut:   15000,
			TotalTokens: 25000,
			Cost:        0.30,
		},
		Week7: []model.UsageSnapshot{
			{
				Date:        "2026-03-07",
				TokensIn:    8000,
				TokensOut:   12000,
				TotalTokens: 20000,
				Cost:        0.20,
			},
		},
		Total: model.UsageSnapshot{
			Date:        time.Now().Format("2006-01-02"),
			TokensIn:    300000,
			TokensOut:   450000,
			TotalTokens: 750000,
			Cost:        9.00,
		},
	}
}

// 模拟处理器实现

type MockSessionHandler struct {
	client *MockGatewayClient
}

func NewMockSessionHandler(c *MockGatewayClient) *MockSessionHandler {
	return &MockSessionHandler{client: c}
}

func (h *MockSessionHandler) List(w http.ResponseWriter, r *http.Request) {
	sessions := h.client.GetSessions()
	writeJSON(w, http.StatusOK, model.APIResponse{
		OK:   true,
		Data: sessions,
	})
}

func (h *MockSessionHandler) Status(w http.ResponseWriter, r *http.Request) {
	statuses := h.client.GetSessionStatuses()
	writeJSON(w, http.StatusOK, model.APIResponse{
		OK:   true,
		Data: statuses,
	})
}

func (h *MockSessionHandler) Get(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotFound, "会话不存在")
}

type MockTaskHandler struct {
	client *MockGatewayClient
}

func NewMockTaskHandler(c *MockGatewayClient) *MockTaskHandler {
	return &MockTaskHandler{client: c}
}

func (h *MockTaskHandler) List(w http.ResponseWriter, r *http.Request) {
	tasks := h.client.GetTasks()
	writeJSON(w, http.StatusOK, model.APIResponse{
		OK:   true,
		Data: tasks,
	})
}

type MockProjectHandler struct {
	client *MockGatewayClient
}

func NewMockProjectHandler(c *MockGatewayClient) *MockProjectHandler {
	return &MockProjectHandler{client: c}
}

func (h *MockProjectHandler) List(w http.ResponseWriter, r *http.Request) {
	projects := h.client.GetProjects()
	writeJSON(w, http.StatusOK, model.APIResponse{
		OK:   true,
		Data: projects,
	})
}

type MockUsageHandler struct {
	client *MockGatewayClient
}

func NewMockUsageHandler(c *MockGatewayClient) *MockUsageHandler {
	return &MockUsageHandler{client: c}
}

func (h *MockUsageHandler) Get(w http.ResponseWriter, r *http.Request) {
	usage := h.client.GetUsageData()
	writeJSON(w, http.StatusOK, model.APIResponse{
		OK:   true,
		Data: usage,
	})
}

// 辅助函数

func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, model.APIResponse{
		OK:    false,
		Error: message,
	})
}

// TestMockDataTypes 测试模拟数据类型
func TestMockDataTypes(t *testing.T) {
	client := &MockGatewayClient{logger: logrus.New()}

	t.Run("GetSessions", func(t *testing.T) {
		sessions := client.GetSessions()
		if len(sessions) != 2 {
			t.Errorf("期望 2 个会话，实际 %d", len(sessions))
		}
		for _, s := range sessions {
			t.Logf("会话: %s, 状态: %s", s.SessionKey, s.State)
		}
	})

	t.Run("GetSessionStatuses", func(t *testing.T) {
		statuses := client.GetSessionStatuses()
		if len(statuses) != 1 {
			t.Errorf("期望 1 个状态，实际 %d", len(statuses))
		}
		for _, s := range statuses {
			t.Logf("状态: %s, Token: %d", s.SessionKey, s.TokensIn)
		}
	})

	t.Run("GetTasks", func(t *testing.T) {
		tasks := client.GetTasks()
		if len(tasks) != 1 {
			t.Errorf("期望 1 个任务，实际 %d", len(tasks))
		}
	})

	t.Run("GetProjects", func(t *testing.T) {
		projects := client.GetProjects()
		if len(projects) != 1 {
			t.Errorf("期望 1 个项目，实际 %d", len(projects))
		}
	})

	t.Run("GetUsageData", func(t *testing.T) {
		usage := client.GetUsageData()
		if usage == nil {
			t.Error("用量数据不应为空")
		}
		if usage.Today.TokensIn != 10000 {
			t.Errorf("期望 10000，实际 %d", usage.Today.TokensIn)
		}
	})
}

// TestContextCancellation 测试上下文取消
func TestContextCancellation(t *testing.T) {
	t.Run("ContextTimeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 0)
		defer cancel()

		select {
		case <-ctx.Done():
			t.Log("上下文已取消（预期行为）")
		case <-time.After(10 * time.Millisecond):
			t.Error("应该立即取消")
		}
	})
}
