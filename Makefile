# Higress Admin SDK Go Makefile
# Copyright (c) 2022-2024 Alibaba Group Holding Ltd.

.PHONY: build test lint clean fmt vet deps all test-coverage help

# Go 参数
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=gofmt
GOLINT=golangci-lint

# 项目名称
PROJECT_NAME=higress-admin-sdk-go

# 默认目标
.DEFAULT_GOAL := help

## build: 构建项目
build:
	@echo ">>> 构建项目..."
	$(GOBUILD) ./...

## test: 运行所有测试
test:
	@echo ">>> 运行测试..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...

## test-coverage: 生成测试覆盖率报告
test-coverage: test
	@echo ">>> 生成覆盖率报告..."
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

## test-coverage-func: 显示测试覆盖率统计
test-coverage-func: test
	@echo ">>> 测试覆盖率统计..."
	$(GOCMD) tool cover -func=coverage.out

## lint: 代码检查
lint:
	@echo ">>> 运行代码检查..."
	$(GOLINT) run ./...

## fmt: 格式化代码
fmt:
	@echo ">>> 格式化代码..."
	$(GOFMT) -w -s .

## vet: 静态检查
vet:
	@echo ">>> 静态检查..."
	$(GOCMD) vet ./...

## clean: 清理构建产物
clean:
	@echo ">>> 清理..."
	$(GOCLEAN)
	rm -f coverage.out coverage.html

## deps: 安装依赖
deps:
	@echo ">>> 安装依赖..."
	$(GOCMD) mod tidy
	$(GOCMD) mod download

## all: 运行所有检查
all: deps fmt vet test
	@echo ">>> 所有检查完成"

## help: 显示帮助信息
help:
	@echo "Higress Admin SDK Go - 构建命令"
	@echo ""
	@echo "使用方法:"
	@echo "  make [target]"
	@echo ""
	@echo "可用目标:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'