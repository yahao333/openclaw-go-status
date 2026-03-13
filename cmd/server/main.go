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

	// 创建 Gateway 客户端
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
	templateDir := "templates"
	if envTemplateDir := os.Getenv("TEMPLATE_DIR"); envTemplateDir != "" {
		templateDir = envTemplateDir
	}
	templateHandler, err := handler.NewTemplateHandler(templateDir)
	if err != nil {
		log.Warnf("加载模板失败: %v，使用纯 JSON API 模式", err)
	}

	// ==================== 路由设置 ====================

	// 创建多路复用器
	mux := http.NewServeMux()

	// 静态文件服务
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// 页面路由（模板）
	if templateHandler != nil {
		mux.HandleFunc("/", templateHandler.Home)
		mux.HandleFunc("/sessions", templateHandler.Sessions)
		mux.HandleFunc("/tasks", templateHandler.Tasks)
		mux.HandleFunc("/projects", templateHandler.Projects)
		mux.HandleFunc("/usage", templateHandler.Usage)
	}

	// API 路由 - 健康检查
	mux.HandleFunc("/health", healthHandler.Check)

	// API 路由 - 会话
	mux.HandleFunc("/api/sessions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			sessionHandler.List(w, r)
		default:
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/sessions/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			sessionHandler.Get(w, r)
		default:
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		}
	})

	// API 路由 - 会话状态
	mux.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			sessionHandler.Status(w, r)
		default:
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		}
	})

	// API 路由 - 任务
	mux.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			taskHandler.List(w, r)
		default:
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		}
	})

	// API 路由 - 项目
	mux.HandleFunc("/api/projects", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			projectHandler.List(w, r)
		default:
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		}
	})

	// API 路由 - 用量
	mux.HandleFunc("/api/usage", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			usageHandler.Get(w, r)
		default:
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		}
	})

	// ==================== 服务器启动 ====================

	// 创建 HTTP 服务器
	server := &http.Server{
		Addr:         cfg.GetAddress(),
		Handler:      mux,
		ReadTimeout:  30 * time.Second,  // 读取超时
		WriteTimeout: 30 * time.Second, // 写入超时
		IdleTimeout:  60 * time.Second, // 空闲超时
	}

	// 启动服务器(在后台)
	go func() {
		log.Infof("服务器启动完成，监听地址: %s", cfg.GetAddress())
		log.Infof("访问 http://localhost:%d 查看 Web 界面", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// ==================== 优雅关闭 ====================

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Infof("收到关闭信号，开始优雅关闭...")

	// 设置关闭超时
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 关闭服务器
	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("服务器关闭失败: %v", err)
	}

	log.Infof("========== OpenClaw Go Status 已关闭 ==========")
}

// getTemplateDir 获取模板目录的绝对路径
func getTemplateDir() string {
	// 获取可执行文件所在目录
	execPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "templates"
	}

	// 尝试多个可能的模板位置
	possiblePaths := []string{
		filepath.Join(execPath, "templates"),
		filepath.Join(execPath, "..", "templates"),
		"templates",
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return "templates"
}
