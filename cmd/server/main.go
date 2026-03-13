// Package main 是程序的入口点
// 负责初始化配置、日志、客户端和 HTTP 服务器
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/yahao333/openclaw-go-status/internal/client"
	"github.com/yahao333/openclaw-go-status/internal/config"
	"github.com/yahao333/openclaw-go-status/internal/handler"
	"github.com/yahao333/openclaw-go-status/internal/logger"
)

func main() {
	// ==================== 初始化阶段 ====================

	// 获取工作目录
	workDir, err := os.Getwd()
	if err != nil {
		workDir = "."
	}
	fmt.Printf("工作目录: %s\n", workDir)

	// 加载配置文件
	configPath := "config.yaml"
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		configPath = envPath
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "加载配置文件失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	log := logger.Init(
		cfg.Logging.Level,
		cfg.Logging.Format,
		cfg.Logging.Output,
		cfg.Logging.File,
	)

	log.Infof("========== OpenClaw Go Status 启动 ==========")
	log.Infof("服务器配置: %s", cfg.GetAddress())
	log.Infof("Gateway 地址: %s", cfg.Gateway.URL)
	log.Infof("日志级别: %s", cfg.Logging.Level)

	// ==================== 组件初始化 ====================

	// 创建 Gateway 客户端（使用 CLI）
	gatewayClient := client.NewGatewayClient(
		cfg.Gateway.URL,
		cfg.Gateway.Timeout,
		log,
	)

	// ==================== 处理器初始化 ====================

	// 创建 API 处理器
	sessionHandler := handler.NewSessionHandler(gatewayClient, log)
	taskHandler := handler.NewTaskHandler(gatewayClient, log)
	projectHandler := handler.NewProjectHandler(gatewayClient, log)
	usageHandler := handler.NewUsageHandler(gatewayClient, log)
	healthHandler := handler.NewHealthHandler(gatewayClient, log)

	// 创建模板处理器
	templateDir := getTemplateDir()
	log.Infof("模板目录: %s", templateDir)

	templateHandler, err := handler.NewTemplateHandler(templateDir, gatewayClient, log)
	if err != nil {
		log.Warnf("加载模板失败: %v，使用纯 JSON API 模式", err)
		templateHandler = nil
	}

	// ==================== 路由设置 ====================

	// 创建多路复用器
	mux := http.NewServeMux()

	// 静态文件服务
	staticDir := filepath.Join(workDir, "static")
	if _, err := os.Stat(staticDir); err == nil {
		fs := http.FileServer(http.Dir(staticDir))
		mux.Handle("/static/", http.StripPrefix("/static/", fs))
	}

	// 页面路由（模板）
	if templateHandler != nil {
		mux.HandleFunc("/", templateHandler.Home)
		mux.HandleFunc("/sessions", templateHandler.Sessions)
		mux.HandleFunc("/tasks", templateHandler.Tasks)
		mux.HandleFunc("/projects", templateHandler.Projects)
		mux.HandleFunc("/usage", templateHandler.Usage)
	}

	// API 路由
	mux.HandleFunc("/health", healthHandler.Check)
	mux.HandleFunc("/api/sessions", sessionHandler.List)
	mux.HandleFunc("/api/sessions/", sessionHandler.Get)
	mux.HandleFunc("/api/status", sessionHandler.Status)
	mux.HandleFunc("/api/tasks", taskHandler.List)
	mux.HandleFunc("/api/projects", projectHandler.List)
	mux.HandleFunc("/api/usage", usageHandler.Get)

	// ==================== 服务器启动 ====================

	server := &http.Server{
		Addr:         cfg.GetAddress(),
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Infof("服务器启动完成，监听地址: %s", cfg.GetAddress())
		log.Infof("访问 http://localhost:%d 查看 Web 界面", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// ==================== 优雅关闭 ====================

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Infof("收到关闭信号，开始优雅关闭...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("服务器关闭失败: %v", err)
	}

	log.Infof("========== OpenClaw Go Status 已关闭 ==========")
}

// getTemplateDir 获取模板目录的绝对路径
func getTemplateDir() string {
	// 先尝试工作目录下的 templates
	workDir, _ := os.Getwd()
	candidates := []string{
		filepath.Join(workDir, "templates"),
		filepath.Join(workDir, "..", "templates"),
		"./templates",
		"templates",
	}

	// 也检查可执行文件目录
	execPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err == nil {
		candidates = append(candidates, filepath.Join(execPath, "templates"))
	}

	for _, path := range candidates {
		absPath, _ := filepath.Abs(path)
		if _, err := os.Stat(absPath); err == nil {
			return absPath
		}
	}

	// 返回默认路径
	return "templates"
}
