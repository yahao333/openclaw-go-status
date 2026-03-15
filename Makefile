.DEFAULT_GOAL := help

# 项目配置
BINARY := openclaw-go-status
CMD := ./cmd/server
FRONTEND_DIR := frontend
EMBEDDED_DIR := internal/handler/frontend_dist

# 颜色定义
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
NC := \033[0m # No Color

.PHONY: help
help:
	@printf "%s\n" \
		"${GREEN}OpenClaw Go Status${NC} - 构建目标:" \
		"" \
		"  ${GREEN}前端相关:${NC}" \
		"    make frontend-deps    安装前端依赖" \
		"    make frontend-build   构建前端 (嵌入到二进制)" \
		"    make frontend-dev     前端开发模式 (Vite)" \
		"    make frontend-clean  清理前端构建产物" \
		"" \
		"  ${GREEN}后端相关:${NC}" \
		"    make build           构建二进制 (包含嵌入的前端)" \
		"    make run             运行服务" \
		"    make dev             开发模式 (前端+后端)" \
		"    make test            运行测试" \
		"    make test-cover      运行测试 (含覆盖率)" \
		"" \
		"  ${GREEN}代码质量:${NC}" \
		"    make fmt             格式化代码" \
		"    make vet             运行 go vet" \
		"    make tidy            整理依赖" \
		"" \
		"  ${GREEN}清理:${NC}" \
		"    make clean           清理所有构建产物" \
		"    make clean-all       清理所有 (包含前端)"

# ==================== 前端相关 ====================

.PHONY: frontend-deps
frontend-deps:
	@echo "${YELLOW}安装前端依赖...${NC}"
	@cd $(FRONTEND_DIR) && npm install

.PHONY: frontend-build
frontend-build: frontend-deps
	@echo "${YELLOW}构建前端...${NC}"
	@cd $(FRONTEND_DIR) && npm run build
	@echo "${GREEN}前端构建完成 -> $(EMBEDDED_DIR)${NC}"

.PHONY: frontend-dev
frontend-dev:
	@echo "${YELLOW}启动前端开发服务器...${NC}"
	@echo "${GREEN}后端 API 代理到 http://localhost:4311${NC}"
	@cd $(FRONTEND_DIR) && npm run dev

.PHONY: frontend-clean
frontend-clean:
	@echo "${YELLOW}清理前端构建产物...${NC}"
	@rm -rf $(FRONTEND_DIR)/dist
	@rm -rf embedded/frontend
	@cd $(FRONTEND_DIR) && rm -rf node_modules 2>/dev/null || true
	@echo "${GREEN}前端清理完成${NC}"

# ==================== 后端相关 ====================

.PHONY: build
build: frontend-build
	@echo "${YELLOW}构建后端...${NC}"
	@go build -o $(BINARY) $(CMD)
	@echo "${GREEN}构建完成: ./$(BINARY)${NC}"

.PHONY: build-only
build-only:
	@echo "${YELLOW}构建后端 (跳过前端)...${NC}"
	@go build -o $(BINARY) $(CMD)
	@echo "${GREEN}构建完成: ./$(BINARY)${NC}"

.PHONY: run
run: build
	@echo "${YELLOW}启动服务...${NC}"
	@./$(BINARY)

.PHONY: dev
dev:
	@echo "${YELLOW}启动开发模式...${NC}"
	@echo "${YELLOW}提示: 前端开发服务器运行在 http://localhost:5173${NC}"
	@echo "${YELLOW}       后端 API 运行在 http://localhost:4311${NC}"
	@go run $(CMD)

.PHONY: test
test:
	@echo "${YELLOW}运行测试...${NC}"
	@go test ./...

.PHONY: test-cover
test-cover:
	@echo "${YELLOW}运行测试 (覆盖率)...${NC}"
	@go test -cover ./...

# ==================== 代码质量 ====================

.PHONY: fmt
fmt:
	@gofmt -w .
	@echo "${GREEN}代码格式化完成${NC}"

.PHONY: fmt-check
fmt-check:
	@if [ -z "$$(gofmt -l .)" ]; then \
		echo "${GREEN}代码格式正确${NC}"; \
	else \
		echo "${RED}代码格式有问题:${NC}"; \
		gofmt -l .; \
		exit 1; \
	fi

.PHONY: vet
vet:
	@go vet ./...

.PHONY: tidy
tidy:
	@go mod tidy
	@echo "${GREEN}依赖整理完成${NC}"

# ==================== 清理 ====================

.PHONY: clean
clean:
	@echo "${YELLOW}清理构建产物...${NC}"
	@rm -f $(BINARY)
	@echo "${GREEN}清理完成${NC}"

.PHONY: clean-all
clean-all: clean
	@echo "${YELLOW}清理所有 (包含前端)...${NC}"
	@rm -rf $(FRONTEND_DIR)/dist
	@rm -rf embedded/frontend
	@cd $(FRONTEND_DIR) && rm -rf node_modules 2>/dev/null || true
	@echo "${GREEN}清理完成${NC}"
