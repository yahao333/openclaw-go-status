# OpenClaw Go Status

OpenClaw 状态查看工具 - 使用 Go 语言重构的轻量级状态监控服务。

## 项目简介

OpenClaw Go Status 是一个专注于状态查看的轻量级服务，提供以下功能：

- ✅ 会话列表查看
- ✅ 会话状态查看（模型、Token 消耗、费用）
- ✅ 任务列表查看
- ✅ 项目列表查看
- ✅ 用量统计查看
- ✅ Cron 任务查看
- ✅ 审批列表查看
- ✅ 异常列表查看
- ✅ 健康检查

## 设计原则

- **只读优先**：仅实现状态查看功能，不执行任何写操作
- **无需鉴权**：无需登录，无需 Token 验证
- **极简设计**：依赖少、性能高、易于部署
- **中文注释**：代码注释使用简体中文

## 技术栈

| 类别 | 技术/工具 |
|------|----------|
| 编程语言 | Go 1.25.0 |
| HTTP 框架 | 标准库 net/http |
| 配置文件 | YAML (gopkg.in/yaml.v3) |
| 日志库 | logrus |
| 测试框架 | testing (标准库) |

## 快速开始

### 1. 克隆项目

```bash
git clone https://github.com/openclaw/openclaw-go-status.git
cd openclaw-go-status
```

### 2. 配置

编辑 `config.yaml` 文件：

```yaml
server:
  host: "0.0.0.0"
  port: 4311

gateway:
  url: "ws://127.0.0.1:18789"
  timeout: 30

logging:
  level: "info"
  format: "json"
  output: "stdout"
```

### 3. 构建

```bash
go build -o openclaw-go-status ./cmd/server
```

### 4. 运行

```bash
./openclaw-go-status
```

服务将在 `http://localhost:4311` 启动。

## API 端点

| 端点 | 方法 | 描述 |
|------|------|------|
| `/health` | GET | 健康检查 |
| `/api/sessions` | GET | 获取会话列表 |
| `/api/sessions/:id` | GET | 获取会话详情 |
| `/api/status` | GET | 获取会话状态 |
| `/api/tasks` | GET | 获取任务列表 |
| `/api/projects` | GET | 获取项目列表 |
| `/api/usage` | GET | 获取用量统计 |

## 日志等级

| 等级 | 说明 |
|------|------|
| DEBUG | 调试信息 |
| INFO | 一般信息 |
| WARN | 警告信息 |
| ERROR | 错误信息 |

## 配置文件说明

### server

- `host`: 监听地址（默认: 0.0.0.0）
- `port`: 监听端口（默认: 4311）

### gateway

- `url`: OpenClaw Gateway WebSocket 地址
- `timeout`: 请求超时时间（秒）

### logging

- `level`: 日志级别（debug, info, warn, error）
- `format`: 日志格式（json, text）
- `output`: 输出位置（stdout, file）
- `file`: 日志文件路径

### polling

- `sessions`: 会话列表轮询间隔（毫秒）
- `status`: 会话状态轮询间隔（毫秒）
- `cron`: Cron 轮询间隔（毫秒）
- `approvals`: 审批轮询间隔（毫秒）

## 运行测试

```bash
go test ./...
```

运行测试并查看覆盖率：

```bash
go test -cover ./...
```

## 项目结构

```
openclaw-go-status/
├── cmd/
│   └── server/
│       └── main.go           # 程序入口
├── internal/
│   ├── config/               # 配置加载
│   │   ├── config.go
│   │   └── config_test.go
│   ├── client/               # Gateway 客户端
│   │   └── gateway.go
│   ├── handler/              # HTTP 处理器
│   │   ├── session.go
│   │   ├── task.go
│   │   ├── project.go
│   │   ├── usage.go
│   │   ├── health.go
│   │   └── helper.go
│   ├── model/                # 数据模型
│   │   ├── types.go
│   │   └── types_test.go
│   └── logger/               # 日志模块
│       ├── logger.go
│       └── logger_test.go
├── config.yaml               # 配置文件
├── go.mod                   # Go 模块文件
├── go.sum                   # 依赖校验
└── README.md                 # 本文件
```

## 许可证

MIT License - 详见 LICENSE 文件

## 注意事项

1. 本项目仅提供状态查看功能，不执行任何写操作
2. 无需登录或鉴权，部署时注意网络安全
3. 建议仅在内网环境使用
4. Gateway 不可用时，会返回模拟数据用于演示

---

*本项目由 OpenClaw 团队维护*
