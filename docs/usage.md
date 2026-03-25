# Higress Admin SDK Go 使用指南

本文档详细介绍如何使用 Higress Admin SDK for Go 管理 Higress 网关配置。

## 目录

- [快速开始](#快速开始)
- [配置说明](#配置说明)
- [服务使用](#服务使用)
- [错误处理](#错误处理)
- [最佳实践](#最佳实践)
- [常见问题](#常见问题)

---

## 快速开始

### 安装

```bash
go get github.com/Jayj1997/higress-admin-sdk-golang/v2
```

### 基本使用

```go
package main

import (
    "context"
    "fmt"
    "log"

    sdk "github.com/Jayj1997/higress-admin-sdk-golang/v2"
    "github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/config"
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

    // 使用服务
    domains, err := provider.DomainService().List(context.Background(), &model.CommonPageQuery{})
    if err != nil {
        log.Fatalf("列出域名失败: %v", err)
    }

    fmt.Printf("找到 %d 个域名\n", domains.Total)
}
```

---

## 配置说明

### Kubernetes 认证方式

SDK 支持多种 Kubernetes 认证方式：

#### 1. kubeconfig 文件路径

```go
cfg := config.NewHigressServiceConfig(
    config.WithKubeConfigPath("/path/to/kubeconfig"),
)
```

#### 2. kubeconfig 内容字符串

```go
kubeconfigContent := `apiVersion: v1
kind: Config
...`

cfg := config.NewHigressServiceConfig(
    config.WithKubeConfigContent(kubeconfigContent),
)
```

#### 3. InCluster ServiceAccount

在 Kubernetes 集群内部运行时，SDK 会自动使用 ServiceAccount：

```go
cfg := config.NewHigressServiceConfig(
    config.WithControllerNamespace("higress-system"),
)
// 不提供 kubeconfig，SDK 会自动使用 InCluster 配置
```

#### 4. Controller 访问 Token

```go
cfg := config.NewHigressServiceConfig(
    config.WithControllerServiceHost("higress-controller.higress-system.svc.cluster.local"),
    config.WithControllerServicePort(15014),
    config.WithControllerAccessToken("your-token"),
)
```

### 控制器配置

| 配置选项 | 说明 | 默认值 |
|---------|------|--------|
| `WithControllerNamespace(ns)` | Higress 控制器命名空间 | `higress-system` |
| `WithControllerWatchedNamespace(ns)` | 控制器监听的命名空间 | `` (所有命名空间) |
| `WithControllerWatchedIngressClassName(name)` | IngressClass 过滤 | `` |
| `WithControllerServiceHost(host)` | 控制器服务主机 | - |
| `WithControllerServicePort(port)` | 控制器服务端口 | `15014` |
| `WithControllerJwtPolicy(policy)` | JWT 策略 | `first-party-jwt` |
| `WithControllerAccessToken(token)` | 访问令牌 | - |

### WASM 插件配置

```go
cfg := config.NewHigressServiceConfig(
    config.WithWasmPluginConfig(&config.WasmPluginServiceConfig{
        ImageRegistryURL:      "oci://your-registry.com/plugins",
        ImagePullSecret:       "your-secret",
        DefaultImagePullPolicy: "IfNotPresent",
    }),
)
```

---

## 服务使用

### 域名管理 (DomainService)

```go
// 列出域名
result, err := provider.DomainService().List(ctx, &model.CommonPageQuery{
    PageNum:  1,
    PageSize: 10,
})

// 获取单个域名
domain, err := provider.DomainService().Get(ctx, "example.com")

// 添加域名
newDomain := &model.Domain{
    Name:        "example.com",
    EnableHTTPS: model.EnableHTTPSOn,
}
domain, err := provider.DomainService().Add(ctx, newDomain)

// 更新域名
domain.EnableHTTPS = model.EnableHTTPSForce
domain, err := provider.DomainService().Update(ctx, domain)

// 删除域名
err = provider.DomainService().Delete(ctx, "example.com")
```

### 路由管理 (RouteService)

```go
// 列出路由
routes, err := provider.RouteService().List(ctx, &model.RoutePageQuery{
    CommonPageQuery: model.CommonPageQuery{
        PageNum:  1,
        PageSize: 10,
    },
    DomainName: "example.com", // 可选：按域名过滤
})

// 获取路由
route, err := provider.RouteService().Get(ctx, "my-route")

// 创建路由
newRoute := &model.Route{
    Name:    "my-route",
    Path:    "/api/v1",
    Domains: []string{"example.com"},
    Services: []model.UpstreamService{
        {Name: "my-service", Namespace: "default", Port: 8080},
    },
}
route, err := provider.RouteService().Add(ctx, newRoute)

// 更新路由
route.Path = "/api/v2"
route, err = provider.RouteService().Update(ctx, route)

// 删除路由
err = provider.RouteService().Delete(ctx, "my-route")
```

#### 路由高级配置

```go
route := &model.Route{
    Name:    "advanced-route",
    Path:    "/api",
    Domains: []string{"example.com"},

    // 路径匹配配置
    PathPrefix: ptrString("/api/v1"),

    // 请求头控制
    HeaderControl: &model.HeaderControlConfig{
        Request: &model.HeaderControlRequest{
            Add:    map[string]string{"X-Custom": "value"},
            Remove: []string{"X-Remove-Me"},
        },
        Response: &model.HeaderControlResponse{
            Add: map[string]string{"X-Response": "value"},
        },
    },

    // CORS 配置
    Cors: &model.CorsConfig{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"GET", "POST"},
        AllowHeaders:     []string{"Content-Type"},
        AllowCredentials: true,
        MaxAge:           3600,
    },

    // 重写配置
    Rewrite: &model.RewriteConfig{
        Path: "/new-path",
    },

    // 重定向配置
    Redirect: &model.RedirectConfig{
        Code:   301,
        Schema: "https",
        Host:   "new.example.com",
    },
}
```

### 服务管理 (ServiceService)

```go
// 列出服务
services, err := provider.ServiceService().List(ctx, &model.CommonPageQuery{
    PageNum:  1,
    PageSize: 20,
})
```

### 服务来源管理 (ServiceSourceService)

```go
// 列出服务来源
sources, err := provider.ServiceSourceService().List(ctx, &model.CommonPageQuery{})

// 添加 Nacos 服务来源
nacosSource := &model.ServiceSource{
    Name:        "my-nacos",
    Type:        model.ServiceSourceTypeNacos2,
    Domain:      "nacos.example.com",
    Port:        8848,
    Namespace:   "public",
    Groups:      []string{"DEFAULT_GROUP"},
    AuthEnabled: true,
    Username:    "nacos",
    Password:    "nacos",
}
source, err := provider.ServiceSourceService().Add(ctx, nacosSource)

// 添加 DNS 服务来源
dnsSource := &model.ServiceSource{
    Name:  "my-dns",
    Type:  model.ServiceSourceTypeDNS,
    Domain: "internal.example.com",
}
source, err := provider.ServiceSourceService().Add(ctx, dnsSource)
```

### TLS 证书管理 (TlsCertificateService)

```go
// 列出证书
certs, err := provider.TlsCertificateService().List(ctx, &model.CommonPageQuery{})

// 添加证书
cert := &model.TlsCertificate{
    Name:        "my-cert",
    Cert:        "-----BEGIN CERTIFICATE-----\n...",
    Key:         "-----BEGIN PRIVATE KEY-----\n...",
    Domains:     []string{"example.com", "www.example.com"},
    ValidBefore: time.Now().Add(365 * 24 * time.Hour),
}
cert, err = provider.TlsCertificateService().Add(ctx, cert)

// 删除证书
err = provider.TlsCertificateService().Delete(ctx, "my-cert")
```

### WASM 插件管理 (WasmPluginService)

```go
// 列出内置插件
plugins, err := provider.WasmPluginService().List(ctx, &model.WasmPluginPageQuery{
    BuiltIn: ptrBool(true),
})

// 获取插件详情
plugin, err := provider.WasmPluginService().Get(ctx, "ai-proxy", "zh-CN")

// 获取插件配置 Schema
config, err := provider.WasmPluginService().GetConfig(ctx, "ai-proxy", "zh-CN")

// 获取插件 README
readme, err := provider.WasmPluginService().GetReadme(ctx, "ai-proxy", "zh-CN")

// 添加自定义插件
customPlugin := &model.WasmPlugin{
    Name:            "my-custom-plugin",
    PluginVersion:   "1.0.0",
    ImageURL:        "oci://my-registry.com/plugins/my-plugin:v1",
    ImagePullPolicy: "IfNotPresent",
}
plugin, err = provider.WasmPluginService().AddCustom(ctx, customPlugin)
```

### WASM 插件实例管理 (WasmPluginInstanceService)

```go
// 创建空实例
instance, err := provider.WasmPluginInstanceService().CreateEmptyInstance(ctx, "ai-proxy")

// 配置全局插件实例
instance.Scope = model.WasmPluginInstanceScopeGlobal
instance.Config = map[string]interface{}{
    "provider": map[string]interface{}{
        "type": "openai",
        "apiTokens": []string{"sk-xxx"},
    },
}
instance, err = provider.WasmPluginInstanceService().AddOrUpdate(ctx, instance)

// 配置域名级插件实例
instance.Scope = model.WasmPluginInstanceScopeDomain
instance.Target = "example.com"
instance, err = provider.WasmPluginInstanceService().AddOrUpdate(ctx, instance)

// 配置路由级插件实例
instance.Scope = model.WasmPluginInstanceScopeRoute
instance.Target = "my-route"
instance, err = provider.WasmPluginInstanceService().AddOrUpdate(ctx, instance)

// 查询插件实例
instance, err = provider.WasmPluginInstanceService().Query(ctx,
    model.WasmPluginInstanceScopeRoute,
    "my-route",
    "ai-proxy",
    nil,
)

// 删除插件实例
err = provider.WasmPluginInstanceService().Delete(ctx,
    model.WasmPluginInstanceScopeRoute,
    "my-route",
    "ai-proxy",
    nil,
)
```

### AI 路由管理 (AiRouteService)

```go
// 创建 AI 路由
aiRoute := &model.AiRoute{
    Name: "my-ai-route",
    Upstreams: []model.AiUpstream{
        {
            ProviderType: "openai",
            ModelMapping: map[string]string{
                "gpt-4": "gpt-4-turbo",
            },
        },
    },
}
route, err := provider.AiRouteService().Add(ctx, aiRoute)
```

### LLM 提供商管理 (LlmProviderService)

```go
// 添加 OpenAI 提供商
provider := &model.LlmProvider{
    Name:         "my-openai",
    Type:         "openai",
    ApiTokens:    []string{"sk-xxx"},
    ModelMapping: map[string]string{"gpt-4": "gpt-4-turbo"},
}
provider, err := provider.LlmProviderService().AddOrUpdate(ctx, provider)

// 添加 Azure OpenAI 提供商
azureProvider := &model.LlmProvider{
    Name:        "my-azure",
    Type:        "azure",
    ApiTokens:   []string{"azure-key"},
    AzureBaseUrl: "https://your-resource.openai.azure.com",
}
provider, err = provider.LlmProviderService().AddOrUpdate(ctx, azureProvider)

// 添加通义千问提供商
qwenProvider := &model.LlmProvider{
    Name:      "my-qwen",
    Type:      "qwen",
    ApiTokens: []string{"qwen-api-key"},
}
provider, err = provider.LlmProviderService().AddOrUpdate(ctx, qwenProvider)
```

### 消费者管理 (ConsumerService)

```go
// 列出消费者
consumers, err := provider.ConsumerService().List(ctx)

// 添加消费者
consumer := &model.Consumer{
    Name: "my-consumer",
    Credentials: []model.Credential{
        &model.KeyAuthCredential{
            Name:   "api-key",
            Source: model.KeyAuthCredentialSourceBearer,
            Key:    "secret-api-key",
        },
    },
}
consumer, err = provider.ConsumerService().AddOrUpdate(ctx, consumer)

// 删除消费者
err = provider.ConsumerService().Delete(ctx, "my-consumer")
```

### MCP 服务器管理 (McpServerService)

```go
// 列出 MCP 服务器
servers, err := provider.McpServerService().List(ctx, &model.McpServerPageQuery{})

// 添加 OpenAPI 类型的 MCP 服务器
mcpServer := &model.McpServer{
    Name:        "my-mcp-server",
    Type:        model.McpServerTypeOpenAPI,
    Description: "My MCP Server",
    Config: map[string]interface{}{
        "openapi_url": "https://api.example.com/openapi.json",
    },
}
server, err := provider.McpServerService().Add(ctx, mcpServer)

// 添加数据库类型的 MCP 服务器
dbServer := &model.McpServer{
    Name:        "my-db-mcp",
    Type:        model.McpServerTypeDatabase,
    Description: "Database MCP Server",
    Config: map[string]interface{}{
        "db_type":     "mysql",
        "host":        "mysql.example.com",
        "port":        3306,
        "database":    "mydb",
        "username":    "user",
        "password":    "pass",
    },
}
server, err = provider.McpServerService().Add(ctx, dbServer)
```

---

## 错误处理

SDK 定义了以下错误类型：

```go
import "github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/errors"

// 业务错误
var ErrBusiness = errors.NewBusinessError("操作失败")

// 资源未找到
var ErrNotFound = errors.NewNotFoundError("资源不存在")

// 资源冲突
var ErrConflict = errors.NewResourceConflictError("资源已存在")

// 验证错误
var ErrValidation = errors.NewValidationError("参数验证失败")
```

### 错误处理示例

```go
domain, err := provider.DomainService().Get(ctx, "nonexistent.com")
if err != nil {
    if errors.IsNotFoundError(err) {
        fmt.Println("域名不存在")
    } else if errors.IsResourceConflictError(err) {
        fmt.Println("资源冲突")
    } else if errors.IsValidationError(err) {
        fmt.Println("参数验证失败:", err.Error())
    } else {
        fmt.Println("其他错误:", err.Error())
    }
    return
}
```

---

## 最佳实践

### 1. 使用 Context 超时控制

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

domains, err := provider.DomainService().List(ctx, &model.CommonPageQuery{})
```

### 2. 分页查询大数据集

```go
pageNum := 1
pageSize := 100

for {
    result, err := provider.RouteService().List(ctx, &model.RoutePageQuery{
        CommonPageQuery: model.CommonPageQuery{
            PageNum:  pageNum,
            PageSize: pageSize,
        },
    })
    if err != nil {
        return err
    }

    // 处理数据
    for _, route := range result.Data {
        // ...
    }

    if result.Total <= pageNum*pageSize {
        break
    }
    pageNum++
}
```

### 3. 优雅关闭

```go
// SDK 内部使用 client-go，会自动处理连接关闭
// 如果需要手动控制，可以使用 context 取消
ctx, cancel := context.WithCancel(context.Background())

// 在程序退出时
cancel()
```

---

## 常见问题

### Q: 如何在 Kubernetes 集群内使用 SDK？

A: 不需要提供 kubeconfig，SDK 会自动使用 InCluster 配置：

```go
cfg := config.NewHigressServiceConfig(
    config.WithControllerNamespace("higress-system"),
)
```

### Q: 如何调试 SDK 请求？

A: 可以通过设置环境变量启用 client-go 日志：

```bash
export KUBE_LOG_LEVEL=10
```

### Q: 支持哪些 LLM 提供商？

A: SDK 支持 20+ 种 LLM 提供商，包括：
- OpenAI
- Azure OpenAI
- 通义千问 (Qwen)
- 文心一言 (ERNIE)
- 智谱 AI (GLM)
- Ollama
- AWS Bedrock
- GCP Vertex AI
- 等等...

### Q: 如何获取 WASM 插件的配置 Schema？

A: 使用 `WasmPluginService.GetConfig` 方法：

```go
config, err := provider.WasmPluginService().GetConfig(ctx, "ai-proxy", "zh-CN")
// config.Config 包含 JSON Schema 格式的配置定义
```

### Q: 如何处理多命名空间？

A: 通过配置 `WithControllerWatchedNamespace` 限制监听范围：

```go
cfg := config.NewHigressServiceConfig(
    config.WithControllerWatchedNamespace("my-namespace"),
)
```

---

## 相关链接

- [API 文档](https://pkg.go.dev/github.com/Jayj1997/higress-admin-sdk-golang/v2)
- [迁移指南](migration.md)
- [示例代码](../examples/)
- [Higress 官方文档](https://higress.io/zh-cn/docs/latest/overview/what-is-higress)
