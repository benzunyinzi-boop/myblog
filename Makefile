# ============================================================
# myblog Makefile
# ============================================================
APP         := myblog-server
PKG         := github.com/yinyin/myblog
CMD_DIR     := ./cmd/server
BIN_DIR     := ./bin
GO          := go
GOFLAGS     := -trimpath
LDFLAGS     := -s -w -X 'main.version=$(shell git describe --tags --always 2>/dev/null || echo dev)'

.PHONY: help tidy fmt lint build run test cover clean migrate-up migrate-down seed-admin

help: ## 显示可用命令
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	  awk 'BEGIN{FS=":.*?## "};{printf "  \033[36m%-14s\033[0m %s\n",$$1,$$2}'

tidy: ## 整理依赖
	$(GO) mod tidy

fmt: ## 格式化
	$(GO) fmt ./...
	@command -v goimports >/dev/null && goimports -w . || echo "[skip] goimports 未安装"

lint: ## 静态检查(需 golangci-lint)
	@command -v golangci-lint >/dev/null || { \
	  echo "[skip] golangci-lint 未安装,执行: brew install golangci-lint"; exit 0; }
	golangci-lint run ./...

build: ## 编译二进制到 bin/
	mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(APP) $(CMD_DIR)
	@echo "==> $(BIN_DIR)/$(APP)"

run: ## 本地运行
	$(GO) run $(CMD_DIR)

test: ## 跑全部测试
	$(GO) test -race -count=1 ./...

cover: ## 覆盖率
	$(GO) test -race -coverprofile=coverage.out ./...
	$(GO) tool cover -func=coverage.out | tail -n 1

clean: ## 清理产物
	rm -rf $(BIN_DIR) coverage.out

# ============================================================
# 数据库:简单用 mysql 客户端跑 migrations/*.sql
# 变量可覆盖:make migrate-up MYSQL_USER=root MYSQL_PWD=xxx MYSQL_DB=myblog
# ============================================================
MYSQL_HOST ?= 127.0.0.1
MYSQL_PORT ?= 3306
MYSQL_USER ?= myblog
MYSQL_PWD  ?= myblog123
MYSQL_DB   ?= myblog

MYSQL_CLI   = mysql -h$(MYSQL_HOST) -P$(MYSQL_PORT) -u$(MYSQL_USER) -p$(MYSQL_PWD)

migrate-up: ## 执行全部 up migrations
	@for f in $$(ls migrations/*.up.sql | sort); do \
	  echo "==> apply $$f"; \
	  $(MYSQL_CLI) $(MYSQL_DB) < $$f || exit 1; \
	done

migrate-down: ## 执行全部 down migrations(倒序)
	@for f in $$(ls migrations/*.down.sql | sort -r); do \
	  echo "==> revert $$f"; \
	  $(MYSQL_CLI) $(MYSQL_DB) < $$f || exit 1; \
	done

seed-admin: ## 创建/重置管理员(PW 指定密码,USERNAME 可覆盖,默认 admin)
	@if [ -z "$(PW)" ]; then echo "usage: make seed-admin PW='yourPwd' [USERNAME=admin]"; exit 2; fi
	$(GO) run ./cmd/seed -username $(or $(USERNAME),admin) -password '$(PW)'
