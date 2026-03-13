// Package config 提供配置管理功能的单元测试
package config

import (
	"os"
	"testing"
)

// TestLoad 测试配置加载功能
func TestLoad(t *testing.T) {
	// 创建临时配置文件
	content := `
server:
  host: "127.0.0.1"
  port: 8080

gateway:
  url: "ws://localhost:18789"
  timeout: 30

logging:
  level: "debug"
  format: "json"
  output: "stdout"

polling:
  sessions: 5000
  status: 2000
  cron: 10000
  approvals: 2000
`
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("写入临时文件失败: %v", err)
	}
	tmpFile.Close()

	// 测试加载配置
	cfg, err := Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("加载配置失败: %v", err)
	}

	// 验证配置值
	if cfg.Server.Host != "127.0.0.1" {
		t.Errorf("期望 host 为 127.0.0.1，实际为 %s", cfg.Server.Host)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("期望 port 为 8080，实际为 %d", cfg.Server.Port)
	}
	if cfg.Gateway.URL != "ws://localhost:18789" {
		t.Errorf("期望 url 为 ws://localhost:18789，实际为 %s", cfg.Gateway.URL)
	}
	if cfg.Gateway.Timeout != 30 {
		t.Errorf("期望 timeout 为 30，实际为 %d", cfg.Gateway.Timeout)
	}
	if cfg.Logging.Level != "debug" {
		t.Errorf("期望 level 为 debug，实际为 %s", cfg.Logging.Level)
	}
	if cfg.Logging.Format != "json" {
		t.Errorf("期望 format 为 json，实际为 %s", cfg.Logging.Format)
	}
	if cfg.Logging.Output != "stdout" {
		t.Errorf("期望 output 为 stdout，实际为 %s", cfg.Logging.Output)
	}
}

// TestLoadDefault 测试默认配置
func TestLoadDefault(t *testing.T) {
	// 创建临时配置文件(使用最小配置)
	content := ""
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("写入临时文件失败: %v", err)
	}
	tmpFile.Close()

	// 测试加载配置
	cfg, err := Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("加载配置失败: %v", err)
	}

	// 验证默认值
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("期望默认 host 为 0.0.0.0，实际为 %s", cfg.Server.Host)
	}
	if cfg.Server.Port != 4311 {
		t.Errorf("期望默认 port 为 4311，实际为 %d", cfg.Server.Port)
	}
	if cfg.Gateway.URL != "ws://127.0.0.1:18789" {
		t.Errorf("期望默认 url 为 ws://127.0.0.1:18789，实际为 %s", cfg.Gateway.URL)
	}
	if cfg.Logging.Level != "info" {
		t.Errorf("期望默认 level 为 info，实际为 %s", cfg.Logging.Level)
	}
}

// TestLoadInvalidPath 测试无效路径
func TestLoadInvalidPath(t *testing.T) {
	_, err := Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Error("期望加载失败，但成功了")
	}
}

// TestGetAddress 测试地址获取
func TestGetAddress(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Host: "127.0.0.1",
			Port: 8080,
		},
	}

	addr := cfg.GetAddress()
	expected := "127.0.0.1:8080"
	if addr != expected {
		t.Errorf("期望地址为 %s，实际为 %s", expected, addr)
	}
}
