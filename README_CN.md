# Higress Admin SDK for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/Jayj1997/higress-admin-sdk-golang.svg)](https://pkg.go.dev/github.com/Jayj1997/higress-admin-sdk-golang)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

[English](README.md) | 简体中文

用于管理 Higress 网关配置的 Go SDK，包括域名、路由、服务、证书和 WASM 插件。

> **注意**：该项目移植自 higress-group 官方的 higress-admin-sdk，详情见：[higress-console](https://github.com/higress-group/higress-console) 和 [higress](https://github.com/alibaba/higress)。

## 功能特性

- **域名管理**：创建、更新、删除和列出网关域名
- **路由管理**：配置路由规则，包括路径匹配、请求头控制等
- **服务管理**：列出和管理后端服务
- **服务来源管理**：配置服务发现来源（Nacos、DNS、静态等）
- **TLS 证书管理**：管理 HTTPS 的 SSL/TLS 证书
- **WASM 插件管理**：配置和管理 WASM 插件
- **AI 路由管理**：配置 AI/LLM 路由规则
- **消费者管理**：管理 API 消费者和凭证
- **MCP 服务器管理**：配置 MCP 服务器实例

## 安装

```bash
go get github.com/Jayj1997/higress-admin-sdk-golang
```

## 快速开始

```go
package main

import (
    "context"
    "fmt"
    "log"

    sdk "github.com/Jayj1997/higress-admin-sdk-golang"
    "github.com/Jayj1997/higress-admin-sdk-golang/pkg/config"
)

func main() {
    // 创建配置
    cfg := config.NewHigressServiceConfig(
        config.WithKubeConfigPath("~/.kube/config"),
        config.WithControllerNamespace("higress-system"),
    )

    // 创建服务提供者
    provider, err := sdk.NewHigressServiceProvider(cfg)
    if err != nil {
        log.Fatalf("创建服务提供者失败: %v", err)
    }

    // 列出所有域名
    domains, err := provider.DomainService().List(context.Background())
    if err != nil {
        log.Fatalf("列出域名失败: %v", err)
    }

    fmt.Printf("找到 %d 个域名\n", len(domains))
    for _, domain := range domains {
        fmt.Printf("- %s\n", domain.Name)
    }
}
```

## 配置选项

| 选项 | 描述 | 默认值 |
|------|------|--------|
| `WithKubeConfigPath(path)` | kubeconfig 文件路径 | - |
| `WithKubeConfigContent(content)` | kubeconfig 内容字符串 | - |
| `WithControllerNamespace(ns)` | Higress 控制器命名空间 | `higress-system` |
| `WithControllerServiceHost(host)` | 控制器服务主机 | - |
| `WithControllerServicePort(port)` | 控制器服务端口 | `15014` |

## API 参考

### 域名服务

```go
// 列出所有域名
domains, err := provider.DomainService().List(ctx)

// 获取特定域名
domain, err := provider.DomainService().Get(ctx, "example.com")

// 添加新域名
newDomain := &model.Domain{
    Name:        "example.com",
    EnableHTTPS: model.EnableHTTPSOn,
}
domain, err := provider.DomainService().Add(ctx, newDomain)

// 更新域名
domain.EnableHTTPS = model.EnableHTTPSForce
domain, err := provider.DomainService().Update(ctx, domain)

// 删除域名
err := provider.DomainService().Delete(ctx, "example.com")
```

### 路由服务

```go
// 分页列出路由
routes, err := provider.RouteService().List(ctx, &model.RoutePageQuery{
    DomainName: "example.com",
})

// 获取特定路由
route, err := provider.RouteService().Get(ctx, "my-route")

// 添加新路由
newRoute := &model.Route{
    Name:    "my-route",
    Path:    "/api/v1",
    Domains: []string{"example.com"},
}
route, err := provider.RouteService().Add(ctx, newRoute)
```

## 文档

- [API 文档](https://pkg.go.dev/github.com/Jayj1997/higress-admin-sdk-golang)
- [Java SDK 迁移指南](docs/migration.md)
- [示例代码](examples/)

## 开发

### 环境要求

- Go 1.21+
- Make

### 运行测试

```bash
# 运行所有测试
make test

# 运行测试并生成覆盖率报告
make test-coverage

# 运行代码检查
make lint
```

### 项目结构

```
higress-admin-sdk-golang/
├── api/v1/              # API 定义
├── pkg/
│   ├── client/          # 客户端实现
│   ├── config/          # 配置
│   ├── model/           # 数据模型
│   ├── service/         # 服务实现
│   ├── constant/        # 常量定义
│   └── errors/          # 错误定义
├── internal/
│   ├── kubernetes/      # Kubernetes 客户端和 CRD
│   └── util/            # 内部工具
├── examples/            # 示例代码
└── test/                # 集成测试
```

## 贡献

欢迎贡献！请阅读[贡献指南](CONTRIBUTING.md)了解详情。

## 许可证

本项目采用 Apache License 2.0 许可证 - 详情见 [LICENSE](LICENSE) 文件。

## 相关项目

- [Higress](https://github.com/alibaba/higress) - 云原生 API 网关
- [Higress Console](https://github.com/higress-group/higress-console) - Higress 控制台