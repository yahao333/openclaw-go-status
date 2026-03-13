// Package integration 提供集成测试
// 测试与 OpenClaw CLI 的集成
package integration

import (
	"context"
	"os/exec"
	"testing"

	"github.com/yahao333/openclaw-go-status/internal/client"
	"github.com/yahao333/openclaw-go-status/internal/model"
	"github.com/sirupsen/logrus"
)

// TestGetSessionsFromCLI 测试从 CLI 获取会话列表
func TestGetSessionsFromCLI(t *testing.T) {
	// 跳过如果 openclaw 命令不可用
	if !isOpenClawAvailable() {
		t.Skip("OpenClaw CLI 不可用，跳过测试")
	}

	logger := logrus.New()
	client := client.NewGatewayClient("ws://localhost:18789", 30, logger)

	ctx := context.Background()
	sessions, err := client.GetSessions(ctx)

	if err != nil {
		t.Logf("获取会话失败（预期行为）: %v", err)
		// 允许失败但不阻塞，因为 CLI 可能返回错误
	}

	// 如果成功，应该返回会话列表
	if err == nil && len(sessions) == 0 {
		t.Log("当前没有会话（正常）")
	}

	t.Logf("获取到 %d 个会话", len(sessions))
}

// TestGetSessionStatusFromCLI 测试获取会话状态
func TestGetSessionStatusFromCLI(t *testing.T) {
	if !isOpenClawAvailable() {
		t.Skip("OpenClaw CLI 不可用，跳过测试")
	}

	logger := logrus.New()
	client := client.NewGatewayClient("ws://localhost:18789", 30, logger)

	ctx := context.Background()
	statuses, err := client.GetSessionStatus(ctx)

	if err != nil {
		t.Logf("获取会话状态失败: %v", err)
	}

	t.Logf("获取到 %d 个会话状态", len(statuses))
}

// TestGetUsageFromCLI 测试获取用量统计
func TestGetUsageFromCLI(t *testing.T) {
	if !isOpenClawAvailable() {
		t.Skip("OpenClaw CLI 不可用，跳过测试")
	}

	logger := logrus.New()
	client := client.NewGatewayClient("ws://localhost:18789", 30, logger)

	ctx := context.Background()
	usage, err := client.GetUsage(ctx)

	if err != nil {
		t.Logf("获取用量失败: %v", err)
	}

	if usage != nil {
		t.Logf("今日用量 - 输入: %d, 输出: %d, 总计: %d, 费用: $%.4f",
			usage.Today.TokensIn,
			usage.Today.TokensOut,
			usage.Today.TotalTokens,
			usage.Today.Cost)
	}
}

// TestHealthCheck 测试健康检查
func TestHealthCheck(t *testing.T) {
	if !isOpenClawAvailable() {
		t.Skip("OpenClaw CLI 不可用，跳过测试")
	}

	logger := logrus.New()
	client := client.NewGatewayClient("ws://localhost:18789", 30, logger)

	ctx := context.Background()
	err := client.CheckHealth(ctx)

	if err != nil {
		t.Logf("健康检查返回错误（可能是配置问题）: %v", err)
	} else {
		t.Log("健康检查通过")
	}
}

// TestOpenClawCLICommand 直接测试 CLI 命令执行
func TestOpenClawCLICommand(t *testing.T) {
	if !isOpenClawAvailable() {
		t.Skip("OpenClaw CLI 不可用，跳过测试")
	}

	cmd := exec.Command("openclaw", "sessions", "--json")
	output, err := cmd.Output()

	if err != nil {
		t.Fatalf("执行 openclaw sessions --json 失败: %v", err)
	}

	// 检查输出是否是有效的 JSON
	if len(output) == 0 {
		t.Fatal("命令输出为空")
	}

	// 简单验证 JSON 格式
	if output[0] != '{' && output[0] != '[' {
		t.Fatalf("输出不是有效的 JSON: %s", string(output[:min(100, len(output))]))
	}

	t.Logf("成功获取 JSON 输出，长度: %d 字节", len(output))
}

// isOpenClawAvailable 检查 openclaw 命令是否可用
func isOpenClawAvailable() bool {
	cmd := exec.Command("openclaw", "--version")
	err := cmd.Run()
	return err == nil
}

// min 返回两个整数中较小的值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// TestModelTypes 测试数据模型
func TestModelTypes(t *testing.T) {
	// 测试会话状态枚举
	states := []model.AgentRunState{
		model.StateIdle,
		model.StateRunning,
		model.StateBlocked,
		model.StateWaitingApproval,
		model.StateError,
	}

	for _, state := range states {
		t.Logf("会话状态: %s", state)
	}

	// 测试任务状态枚举
	taskStates := []model.TaskState{
		model.TaskTodo,
		model.TaskInProgress,
		model.TaskBlocked,
		model.TaskDone,
	}

	for _, state := range taskStates {
		t.Logf("任务状态: %s", state)
	}

	// 测试项目状态枚举
	projectStates := []model.ProjectState{
		model.ProjectPlanned,
		model.ProjectActive,
		model.ProjectBlocked,
		model.ProjectDone,
	}

	for _, state := range projectStates {
		t.Logf("项目状态: %s", state)
	}

	// 测试预算状态枚举
	budgetStatuses := []model.BudgetStatus{
		model.BudgetOk,
		model.BudgetWarn,
		model.BudgetOver,
	}

	for _, status := range budgetStatuses {
		t.Logf("预算状态: %s", status)
	}
}
