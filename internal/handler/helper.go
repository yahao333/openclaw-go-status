// Package handler 提供 HTTP 请求处理功能
// 包含通用的辅助函数
package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/openclaw/openclaw-go-status/internal/model"
)

// writeJSON 写入 JSON 响应
// 参数:
//   - w: HTTP 响应写入器
//   - statusCode: HTTP 状态码
//   - data: 响应数据
func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		// 如果编码失败，写入原始错误信息
		http.Error(w, `{"error":"encoding error"}`, http.StatusInternalServerError)
	}
}

// writeError 写入错误响应
// 参数:
//   - w: HTTP 响应写入器
//   - statusCode: HTTP 状态码
//   - message: 错误消息
func writeError(w http.ResponseWriter, statusCode int, message string) {
	response := model.APIResponse{
		OK:    false,
		Error: message,
	}
	writeJSON(w, statusCode, response)
}

// extractID 从 URL 路径中提取 ID
// 例如: /api/sessions/abc123 -> abc123
func extractID(r *http.Request) string {
	// 获取路径
	path := r.URL.Path

	// 按 / 分割
	parts := strings.Split(path, "/")

	// 返回最后一部分
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return ""
}
