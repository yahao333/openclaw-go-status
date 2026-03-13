// Package config 提供配置管理功能
// 负责加载和解析 config.yaml 配置文件
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `yaml:"host"` // 监听地址
	Port int    `yaml:"port"` // 监听端口
}

// GatewayConfig Gateway 配置
type GatewayConfig struct {
	URL     string `yaml:"url"`     // Gateway WebSocket 地址
	Timeout int    `yaml:"timeout"` // 请求超时时间(秒)
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level  string `yaml:"level"`  // 日志级别: debug, info, warn, error
	Format string `yaml:"format"` // 日志格式: json, text
	Output string `yaml:"output"` // 输出位置: stdout, file
	File   string `yaml:"file"`   // 日志文件路径
}

// PollingConfig 轮询配置
type PollingConfig struct {
	Sessions  int `yaml:"sessions"`  // 会话列表轮询间隔(毫秒)
	Status    int `yaml:"status"`    // 会话状态轮询间隔(毫秒)
	Cron      int `yaml:"cron"`      // Cron 轮询间隔(毫秒)
	Approvals int `yaml:"approvals"` // 审批轮询间隔(毫秒)
}

// Config 应用程序配置
type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Gateway GatewayConfig `yaml:"gateway"`
	Logging LoggingConfig `yaml:"logging"`
	Polling PollingConfig `yaml:"polling"`
}

// Load 加载配置文件
// 参数: configPath 配置文件路径
// 返回: *Config 配置指针, error 错误信息
func Load(configPath string) (*Config, error) {
	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析 YAML
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 设置默认值
	if cfg.Server.Host == "" {
		cfg.Server.Host = "0.0.0.0"
	}
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 4311
	}
	if cfg.Gateway.URL == "" {
		cfg.Gateway.URL = "ws://127.0.0.1:18789"
	}
	if cfg.Gateway.Timeout == 0 {
		cfg.Gateway.Timeout = 30
	}
	if cfg.Logging.Level == "" {
		cfg.Logging.Level = "info"
	}
	if cfg.Logging.Format == "" {
		cfg.Logging.Format = "json"
	}
	if cfg.Logging.Output == "" {
		cfg.Logging.Output = "stdout"
	}

	// 设置默认轮询间隔
	if cfg.Polling.Sessions == 0 {
		cfg.Polling.Sessions = 5000
	}
	if cfg.Polling.Status == 0 {
		cfg.Polling.Status = 2000
	}
	if cfg.Polling.Cron == 0 {
		cfg.Polling.Cron = 10000
	}
	if cfg.Polling.Approvals == 0 {
		cfg.Polling.Approvals = 2000
	}

	return &cfg, nil
}

// GetAddress 获取服务器监听地址
// 返回: 地址字符串 (host:port)
func (c *Config) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
