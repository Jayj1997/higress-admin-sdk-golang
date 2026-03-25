# Higress Admin SDK for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/Jayj1997/higress-admin-sdk-golang.svg)](https://pkg.go.dev/github.com/Jayj1997/higress-admin-sdk-golang)
[![Version](https://img.shields.io/badge/version-v2.2.0-blue.svg)](https://github.com/Jayj1997/higress-admin-sdk-golang/releases)

[English](README_EN.md) | 简体中文

用于管理 Higress 网关配置的 Go SDK，包括域名、路由、服务、证书和 WASM 插件。

> **版本说明**：该项目移植自 higress-group 官方的 higress-admin-sdk，版本 **v2.2.0** 与 [higress-console 2.2.0](https://github.com/higress-group/higress-console) 对齐。

## 功能特性

- **域名管理**：创建、更新、删除和列出网关域名
- **路由管理**：配置路由规则，包括路径匹配、请求头控制、CORS、重写等
- **服务管理**：列出和管理后端服务
- **服务来源管理**：配置服务发现来源（Nacos、DNS、静态等）
- **TLS 证书管理**：管理 HTTPS 的 SSL/TLS 证书
- **WASM 插件管理**：配置和管理 WASM 插件（40+ 内置插件）
- **AI 路由管理**：配置 AI/LLM 路由规则（支持 20+ LLM 提供商）
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
    "github.com/Jayj1997/higress-admin-sdk-golang/pkg/model"
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
    domains, err := provider.DomainService().List(context.Background(), &model.CommonPageQuery{})
    if err != nil {
        log.Fatalf("列出域名失败: %v", err)
    }

    fmt.Printf("找到 %d 个域名\n", domains.Total)
    for _, domain := range domains.Data {
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

## 文档

- **[使用指南](docs/usage.md)** - 详细的配置说明和服务使用方法
- **[迁移指南](docs/migration.md)** - Java SDK 到 Go SDK 的迁移说明
- **[API 文档](https://pkg.go.dev/github.com/Jayj1997/higress-admin-sdk-golang)** - GoDoc 生成的 API 参考
- **[变更日志](CHANGELOG.md)** - 版本变更记录

## 示例代码

| 示例 | 说明 |
|------|------|
| [hello-sdk](examples/hello-sdk/) | 基础使用示例 |
| [route-management](examples/route-management/) | 路由管理示例（CRUD、CORS、请求头控制） |
| [ai-route](examples/ai-route/) | AI 路由示例（LLM 提供商管理） |
| [wasm-plugin](examples/wasm-plugin/) | WASM 插件配置示例 |

## API 参考

### 域名服务

```go
// 列出所有域名
domains, err := provider.DomainService().List(ctx, &model.CommonPageQuery{})

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
    Domains: []string{"example.com"},
    Path: &route.RoutePredicate{
        MatchType: route.MatchTypePrefix,
        Path:      "/api/v1",
    },
    Services: []*route.UpstreamService{
        {Name: "my-service", Namespace: "default", Port: 8080},
    },
}
route, err := provider.RouteService().Add(ctx, newRoute)
```

### WASM 插件服务

```go
// 列出内置插件
builtIn := true
plugins, err := provider.WasmPluginService().List(ctx, &model.WasmPluginPageQuery{
    BuiltIn: &builtIn,
})

// 获取插件配置
config, err := provider.WasmPluginService().GetConfig(ctx, "ai-proxy", "zh-CN")

// 创建插件实例
instance, err := provider.WasmPluginInstanceService().CreateEmptyInstance(ctx, "ai-proxy")
instance.Scope = model.WasmPluginInstanceScopeRoute
instance.Target = "my-route"
instance.Configurations = map[string]interface{}{
    "provider": map[string]interface{}{
        "type":      "openai",
        "apiTokens": []string{"sk-xxx"},
    },
}
instance, err = provider.WasmPluginInstanceService().AddOrUpdate(ctx, instance)
```

## 开发

### 环境要求

- Go 1.21+
- Make

### 运行测试

```bash
# 运行单元测试（无需 K8s 集群）
make test

# 运行单元测试（make test 的别名）
make test-unit

# 运行集成测试（需要 K8s 集群）
make test-integration

# 运行所有测试，包括集成测试（需要 K8s 集群）
make test-all

# 运行测试并生成覆盖率报告
make test-coverage

# 运行代码检查
make lint
```

#### 测试类型

| 命令 | 描述 | 需要 K8s |
|------|------|----------|
| `make test` | 仅单元测试 | 否 |
| `make test-unit` | 同 `make test` | 否 |
| `make test-integration` | 仅集成测试 | 是 |
| `make test-all` | 所有测试（单元 + 集成） | 是 |

### 项目结构

```
higress-admin-sdk-golang/
├── CHANGELOG.md           # 变更日志
├── docs/
│   ├── usage.md           # 使用指南
│   └── migration.md       # 迁移指南
├── examples/
│   ├── hello-sdk/         # 基础示例
│   ├── route-management/  # 路由管理示例
│   ├── ai-route/          # AI 路由示例
│   └── wasm-plugin/       # WASM 插件示例
├── pkg/
│   ├── client/            # 客户端实现
│   ├── config/            # 配置
│   ├── model/             # 数据模型
│   ├── service/           # 服务实现
│   ├── constant/          # 常量定义
│   └── errors/            # 错误定义
├── internal/
│   ├── kubernetes/        # Kubernetes 客户端和 CRD
│   └── resources/         # 内置资源（WASM 插件配置）
└── test/                  # 集成测试
```

## 贡献

欢迎贡献！请阅读[贡献指南](CONTRIBUTING.md)了解详情。

## 许可证

本项目采用 Apache License 2.0 许可证 - 详情见 [LICENSE](LICENSE) 文件。

## 相关项目

- [Higress](https://github.com/alibaba/higress) - 云原生 API 网关
- [Higress Console](https://github.com/higress-group/higress-console) - Higress 控制台
