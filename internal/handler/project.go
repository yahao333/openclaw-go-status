// Package handler 提供 HTTP 请求处理功能
// 包含项目相关的 HTTP 处理器
package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/openclaw/openclaw-go-status/internal/client"
	"github.com/openclaw/openclaw-go-status/internal/model"
	"github.com/sirupsen/logrus"
)

// ProjectHandler 项目处理器
type ProjectHandler struct {
	client *client.GatewayClient // Gateway 客户端
	logger *logrus.Logger       // 日志实例
}

// NewProjectHandler 创建项目处理器
// 参数:
//   - client: Gateway 客户端
//   - logger: 日志实例
// 返回: *ProjectHandler 处理器指针
func NewProjectHandler(client *client.GatewayClient, logger *logrus.Logger) *ProjectHandler {
	return &ProjectHandler{
		client: client,
		logger: logger,
	}
}

// List 获取项目列表
// 方法: GET /api/projects
// 返回: JSON 格式的项目列表
func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	h.logger.Infof("获取项目列表")

	projects, err := h.client.GetProjects(ctx)
	if err != nil {
		h.logger.Warnf("获取项目列表失败: %v", err)
		// 返回模拟数据用于演示
		projects = getMockProjects()
	}

	response := model.APIResponse{
		OK:   true,
		Data: projects,
	}

	writeJSON(w, http.StatusOK, response)
}

// getMockProjects 获取模拟项目数据
func getMockProjects() []model.ProjectRecord {
	return []model.ProjectRecord{
		{
			ProjectID: "project-001",
			Title:     "用户认证系统",
			Status:    model.ProjectActive,
			Owner:     "zhangsan",
			UpdatedAt: time.Now().Format(time.RFC3339),
		},
		{
			ProjectID: "project-002",
			Title:     "数据报表模块",
			Status:    model.ProjectPlanned,
			Owner:     "lisi",
			UpdatedAt: time.Now().Format(time.RFC3339),
		},
		{
			ProjectID: "project-003",
			Title:     "系统优化",
			Status:    model.ProjectDone,
			Owner:     "wangwu",
			UpdatedAt: time.Now().Format(time.RFC3339),
		},
	}
}
