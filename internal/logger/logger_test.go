// Package logger 提供日志功能的单元测试
package logger

import (
	"testing"
)

// TestInit 测试日志初始化
func TestInit(t *testing.T) {
	// 测试不同日志级别
	tests := []struct {
		level  string
		format string
		output string
	}{
		{"debug", "json", "stdout"},
		{"info", "text", "stdout"},
		{"warn", "json", "stdout"},
		{"error", "text", "stdout"},
		{"invalid", "json", "stdout"}, // 无效级别应回退到 info
	}

	for _, tt := range tests {
		logger := Init(tt.level, tt.format, tt.output, "")
		if logger == nil {
			t.Errorf("日志初始化失败")
		}
		if logger.Level == 0 {
			t.Errorf("日志级别未设置")
		}
	}
}

// TestWithModule 测试模块日志
func TestWithModule(t *testing.T) {
	logger := Init("info", "text", "stdout", "")
	
	// 测试添加模块
	entry := WithModule(logger, "test-module")
	if entry == nil {
		t.Errorf("创建模块日志失败")
	}
}
