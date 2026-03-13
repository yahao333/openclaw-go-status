// Package logger 提供日志管理功能
// 支持多种日志级别、格式和输出方式
package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

// LevelMap 日志级别映射表
var LevelMap = map[string]logrus.Level{
	"debug": logrus.DebugLevel,
	"info":  logrus.InfoLevel,
	"warn":  logrus.WarnLevel,
	"error": logrus.ErrorLevel,
}

// Init 初始化日志系统
// 参数:
//   - level: 日志级别字符串 (debug, info, warn, error)
//   - format: 日志格式 (json, text)
//   - output: 输出位置 (stdout, file)
//   - filePath: 日志文件路径
// 返回: *logrus.Logger 日志实例
func Init(level, format, output, filePath string) *logrus.Logger {
	logger := logrus.New()

	// 设置日志级别
	lvl, ok := LevelMap[strings.ToLower(level)]
	if !ok {
		lvl = logrus.InfoLevel
	}
	logger.SetLevel(lvl)

	// 设置日志格式
	switch strings.ToLower(format) {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			// 时间戳格式
			TimestampFormat: "2006-01-02T15:04:05Z07:00",
		})
	default:
		logger.SetFormatter(&logrus.TextFormatter{
			// 时间戳格式
			TimestampFormat: "2006-01-02 15:04:05",
			// 启用颜色
			ForceColors: true,
			// 完整时间戳
			FullTimestamp: true,
		})
	}

	// 设置输出目标
	switch strings.ToLower(output) {
	case "file":
		// 创建日志目录
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "创建日志目录失败: %v\n", err)
			logger.SetOutput(os.Stdout)
			return logger
		}

		// 打开日志文件
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
		if err != nil {
			fmt.Fprintf(os.Stderr, "打开日志文件失败: %v\n", err)
			logger.SetOutput(os.Stdout)
			return logger
		}

		// 同时输出到文件和控制台
		logger.SetOutput(io.MultiWriter(os.Stdout, file))
	default:
		logger.SetOutput(os.Stdout)
	}

	return logger
}

// WithModule 创建带有模块名的日志实例
// 参数:
//   - logger: 基础日志实例
//   - module: 模块名称
// 返回: *logrus.Entry 带模块的日志条目
func WithModule(logger *logrus.Logger, module string) *logrus.Entry {
	return logger.WithField("module", module)
}
