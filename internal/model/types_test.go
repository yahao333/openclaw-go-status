// Package model 提供数据模型的单元测试
package model

import (
	"encoding/json"
	"testing"
)

// TestAgentRunState 测试 Agent 运行状态
func TestAgentRunState(t *testing.T) {
	tests := []struct {
		state   AgentRunState
		want    string
	}{
		{StateIdle, "idle"},
		{StateRunning, "running"},
		{StateBlocked, "blocked"},
		{StateWaitingApproval, "waiting_approval"},
		{StateError, "error"},
	}

	for _, tt := range tests {
		got := string(tt.state)
		if got != tt.want {
			t.Errorf("AgentRunState = %v, want %v", got, tt.want)
		}
	}
}

// TestTaskState 测试任务状态
func TestTaskState(t *testing.T) {
	tests := []struct {
		state TaskState
		want  string
	}{
		{TaskTodo, "todo"},
		{TaskInProgress, "in_progress"},
		{TaskBlocked, "blocked"},
		{TaskDone, "done"},
	}

	for _, tt := range tests {
		got := string(tt.state)
		if got != tt.want {
			t.Errorf("TaskState = %v, want %v", got, tt.want)
		}
	}
}

// TestProjectState 测试项目状态
func TestProjectState(t *testing.T) {
	tests := []struct {
		state ProjectState
		want  string
	}{
		{ProjectPlanned, "planned"},
		{ProjectActive, "active"},
		{ProjectBlocked, "blocked"},
		{ProjectDone, "done"},
	}

	for _, tt := range tests {
		got := string(tt.state)
		if got != tt.want {
			t.Errorf("ProjectState = %v, want %v", got, tt.want)
		}
	}
}

// TestBudgetStatus 测试预算状态
func TestBudgetStatus(t *testing.T) {
	tests := []struct {
		status BudgetStatus
		want   string
	}{
		{BudgetOk, "ok"},
		{BudgetWarn, "warn"},
		{BudgetOver, "over"},
	}

	for _, tt := range tests {
		got := string(tt.status)
		if got != tt.want {
			t.Errorf("BudgetStatus = %v, want %v", got, tt.want)
		}
	}
}

// TestSessionSummaryJSON 测试会话摘要 JSON 序列化
func TestSessionSummaryJSON(t *testing.T) {
	session := SessionSummary{
		SessionKey:   "test-001",
		Label:        "测试会话",
		AgentID:      "agent-001",
		State:        StateRunning,
		LastMessageAt: "2026-03-13T10:00:00Z",
	}

	data, err := json.Marshal(session)
	if err != nil {
		t.Fatalf("JSON 序列化失败: %v", err)
	}

	var decoded SessionSummary
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("JSON 反序列化失败: %v", err)
	}

	if decoded.SessionKey != session.SessionKey {
		t.Errorf("SessionKey 不匹配: got %v, want %v", decoded.SessionKey, session.SessionKey)
	}
	if decoded.State != session.State {
		t.Errorf("State 不匹配: got %v, want %v", decoded.State, session.State)
	}
}

// TestHealthResponseJSON 测试健康检查响应 JSON
func TestHealthResponseJSON(t *testing.T) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: "2026-03-13T10:00:00Z",
		Checks: map[string]Check{
			"server": {
				Status:  "pass",
				Message: "OK",
			},
		},
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("JSON 序列化失败: %v", err)
	}

	var decoded HealthResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("JSON 反序列化失败: %v", err)
	}

	if decoded.Status != "healthy" {
		t.Errorf("Status 不匹配")
	}
	if len(decoded.Checks) != 1 {
		t.Errorf("Checks 长度不匹配")
	}
}
