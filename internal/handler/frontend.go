package handler

import (
	"io/fs"
	"net/http"
	"path"
	"strings"

	"embed"
)

//go:embed frontend_dist
var frontendFiles embed.FS

// GetFrontendHandler 返回前端处理器
// 使用嵌入的文件系统，无需外部依赖
func GetFrontendHandler() (http.Handler, error) {
	return &embeddedFrontendHandler{fs: frontendFiles}, nil
}

type embeddedFrontendHandler struct {
	fs embed.FS
}

func (h *embeddedFrontendHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqPath := path.Clean(r.URL.Path)

	// 防止路径遍历
	if strings.HasPrefix(reqPath, "..") {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// 根路径返回 index.html
	if reqPath == "/" {
		reqPath = "/index.html"
	}

	// 添加前缀以匹配嵌入的路径
	filePath := "frontend_dist" + reqPath

	// 尝试打开文件
	data, err := h.fs.ReadFile(filePath)
	if err != nil {
		// SPA fallback: 返回 index.html
		indexData, err := h.fs.ReadFile("frontend_dist/index.html")
		if err != nil {
			http.Error(w, "Frontend not built", http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		w.Write(indexData)
		return
	}

	// 设置 Content-Type 和缓存
	contentType := getContentType(reqPath)
	w.Header().Set("Content-Type", contentType)

	// 静态资源缓存 1 小时，HTML 不缓存
	if !strings.HasSuffix(reqPath, ".html") {
		w.Header().Set("Cache-Control", "public, max-age=3600")
	} else {
		w.Header().Set("Cache-Control", "no-cache")
	}

	// 写入文件内容
	w.Write(data)
}

// getFileInfo 获取嵌入文件的信息
func getFileInfo(fsEmbed embed.FS, name string) (fs.FileInfo, error) {
	file, err := fsEmbed.Open(name)
	if err != nil {
		return nil, err
	}
	return file.Stat()
}

func getContentType(filename string) string {
	ext := strings.ToLower(path.Ext(filename))
	switch ext {
	case ".html":
		return "text/html; charset=utf-8"
	case ".js", ".mjs":
		return "application/javascript"
	case ".css":
		return "text/css"
	case ".json":
		return "application/json"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	case ".ico":
		return "image/x-icon"
	case ".woff":
		return "font/woff"
	case ".woff2":
		return "font/woff2"
	case ".ttf":
		return "font/ttf"
	default:
		return "application/octet-stream"
	}
}
