// Package main 是程序的入口点
// 负责初始化配置、日志、客户端和 HTTP 服务器
package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
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

	// 尝试获取嵌入的前端处理器（默认启用）
	frontendHandler, err := handler.GetFrontendHandler()
	if err != nil {
		log.Warnf("前端资源未找到: %v，使用模板模式", err)
		frontendHandler = nil
	} else {
		log.Infof("启用嵌入的前端模式")
	}

	// API 路由（必须放在前端 handler 之前）
	mux.HandleFunc("/health", healthHandler.Check)
	mux.HandleFunc("/api/sessions", sessionHandler.List)
	mux.HandleFunc("/api/sessions/", sessionHandler.Get)
	mux.HandleFunc("/api/status", sessionHandler.Status)
	mux.HandleFunc("/api/dashboard", sessionHandler.DashboardStats)
	mux.HandleFunc("/api/tasks", taskHandler.List)
	mux.HandleFunc("/api/projects", projectHandler.List)
	mux.HandleFunc("/api/usage", usageHandler.Get)

	// 根据模式选择前端处理
	if frontendHandler != nil {
		// React 前端模式（SPA fallback）
		mux.Handle("/", frontendHandler)
	} else if templateHandler != nil {
		// 模板模式（保留原有逻辑）
		// 静态文件服务
		staticDir := filepath.Join(workDir, "static")
		if _, err := os.Stat(staticDir); err == nil {
			fs := http.FileServer(http.Dir(staticDir))
			mux.Handle("/static/", http.StripPrefix("/static/", fs))
		}

		// 页面路由
		mux.HandleFunc("/", templateHandler.Home)
		mux.HandleFunc("/sessions", templateHandler.Sessions)
		mux.HandleFunc("/tasks", templateHandler.Tasks)
		mux.HandleFunc("/projects", templateHandler.Projects)
		mux.HandleFunc("/usage", templateHandler.Usage)

		// SPA API 路由
		mux.HandleFunc("/api/dashboard", templateHandler.DashboardAPI)
		mux.HandleFunc("/api/sessions-page", templateHandler.SessionsPageAPI)
		mux.HandleFunc("/api/tasks-page", templateHandler.TasksPageAPI)
		mux.HandleFunc("/api/projects-page", templateHandler.ProjectsPageAPI)
		mux.HandleFunc("/api/usage-page", templateHandler.UsagePageAPI)
	} else {
		// 纯 API 模式
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"message": "OpenClaw Status API", "endpoints": ["/api/sessions", "/api/tasks", "/api/projects", "/api/usage"]}`))
		})
	}

	// ==================== 服务器启动 ====================

	server := &http.Server{
		Addr:         cfg.GetAddress(),
		Handler:      withAccessLog(log, mux),
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

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
	bytes      int
}

func (w *statusRecorder) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusRecorder) Write(p []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(p)
	w.bytes += n
	return n, err
}

func withAccessLog(log *logrus.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &statusRecorder{ResponseWriter: w}

		next.ServeHTTP(recorder, r)
		if recorder.statusCode == 0 {
			recorder.statusCode = http.StatusOK
		}

		duration := time.Since(start)

		fields := logrus.Fields{
			"method":      r.Method,
			"path":        r.URL.Path,
			"query":       r.URL.RawQuery,
			"status":      recorder.statusCode,
			"bytes":       recorder.bytes,
			"duration_ms": duration.Milliseconds(),
			"remote_ip":   clientIP(r),
			"referer":     r.Referer(),
			"user_agent":  r.UserAgent(),
		}

		if requestID := r.Header.Get("X-Request-Id"); requestID != "" {
			fields["request_id"] = requestID
		} else if requestID := r.Header.Get("X-Request-ID"); requestID != "" {
			fields["request_id"] = requestID
		}

		log.WithFields(fields).Info("http_access")
	})
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
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
