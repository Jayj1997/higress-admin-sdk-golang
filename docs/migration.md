# Java SDK 到 Go SDK 迁移指南

本文档帮助用户从 higress-admin-sdk Java 版本迁移到 Go 版本。

## 目录

- [概述](#概述)
- [包结构对照](#包结构对照)
- [API 差异对照表](#api-差异对照表)
- [代码示例对照](#代码示例对照)
- [类型映射](#类型映射)
- [异常处理差异](#异常处理差异)

---

## 概述

### 项目背景

higress-admin-sdk-golang 是 higress-admin-sdk Java 版本的 Go 移植版本，与官方 higress&higress-console 2.2.0 版本对齐。

### 主要差异

| 特性 | Java SDK | Go SDK |
|------|----------|--------|
| 语言版本 | Java 8+ | Go 1.21+ |
| 构建工具 | Maven | Go Modules |
| 异常处理 | try-catch | error 返回值 |
| 空值处理 | null | nil |
| 集合类型 | List, Map | slice, map |
| 异步操作 | CompletableFuture | goroutine + channel |

---

## 包结构对照

```
Java                                    Go
---                                     ---
com.alibaba.higress.sdk               → github.com/Jayj1997/higress-admin-sdk-golang

├── config                             → pkg/config
│   └── HigressServiceConfig           → HigressServiceConfig
│
├── constant                           → pkg/constant
│   ├── HigressConstants               → higress.go
│   └── KubernetesConstants            → kubernetes.go
│
├── exception                          → pkg/errors
│   ├── BusinessException              → errors.go
│   ├── NotFoundException              → errors.go
│   ├── ResourceConflictException      → errors.go
│   └── ValidationException            → errors.go
│
├── model                              → pkg/model
│   ├── Domain                         → domain.go
│   ├── Route                          → route.go
│   ├── Service                        → service.go
│   ├── ServiceSource                  → service_source.go
│   ├── TlsCertificate                 → tls_certificate.go
│   ├── WasmPlugin                     → wasm_plugin.go
│   ├── ai/                            → ai_route.go, llm_provider.go
│   ├── consumer/                      → consumer.go, allow_list.go
│   ├── mcp/                           → mcp_server.go
│   └── route/                         → route/route_predicate.go
│
├── service                            → pkg/service
│   ├── HigressServiceProvider         → pkg/client/provider.go
│   ├── DomainService                  → domain_service.go
│   ├── RouteService                   → route_service.go
│   ├── ServiceService                 → service_service.go
│   ├── ServiceSourceService           → service_source_service.go
│   ├── TlsCertificateService          → tls_certificate_service.go
│   ├── WasmPluginService              → wasm_plugin_service.go
│   ├── WasmPluginInstanceService      → wasm_plugin_instance_service.go
│   ├── ai/                            → ai/, llm_provider_service.go
│   ├── consumer/                      → consumer/, consumer_service.go
│   ├── kubernetes/                    → internal/kubernetes/
│   └── mcp/                           → mcp/, mcp_service.go
│
└── util                               → internal/util
```

---

## API 差异对照表

### 配置初始化

**Java:**

```java
import com.alibaba.higress.sdk.config.HigressServiceConfig;
import com.alibaba.higress.sdk.service.HigressServiceProvider;

HigressServiceConfig config = HigressServiceConfig.builder()
    .kubeConfigPath("/path/to/kubeconfig")
    .controllerNamespace("higress-system")
    .build();

HigressServiceProvider provider = HigressServiceProvider.create(config);
```

**Go:**

```go
import (
    sdk "github.com/Jayj1997/higress-admin-sdk-golang"
    "github.com/Jayj1997/higress-admin-sdk-golang/pkg/config"
)

cfg := config.NewHigressServiceConfig(
    config.WithKubeConfigPath("/path/to/kubeconfig"),
    config.WithControllerNamespace("higress-system"),
)

provider, err := sdk.NewHigressServiceProvider(cfg)
if err != nil {
    log.Fatal(err)
}
```

### 服务获取方式

**Java:**

```java
DomainService domainService = provider.domainService();
RouteService routeService = provider.routeService();
```

**Go:**

```go
domainService := provider.DomainService()
routeService := provider.RouteService()
```

### 方法命名差异

| Java 方法 | Go 方法 | 说明 |
|-----------|---------|------|
| `list()` | `List(ctx, query)` | Go 需要传入 context 和查询参数 |
| `get(name)` | `Get(ctx, name)` | Go 需要传入 context |
| `add(entity)` | `Add(ctx, entity)` | Go 需要传入 context |
| `update(entity)` | `Update(ctx, entity)` | Go 需要传入 context |
| `delete(name)` | `Delete(ctx, name)` | Go 需要传入 context |

---

## 代码示例对照

### 域名管理

**Java:**

```java
// 列出域名
List<Domain> domains = provider.domainService().list();

// 获取域名
Domain domain = provider.domainService().get("example.com");

// 添加域名
Domain newDomain = new Domain();
newDomain.setName("example.com");
newDomain.setEnableHTTPS(EnableHTTPS.ON);
domain = provider.domainService().add(newDomain);

// 更新域名
domain.setEnableHTTPS(EnableHTTPS.FORCE);
domain = provider.domainService().update(domain);

// 删除域名
provider.domainService().delete("example.com");
```

**Go:**

```go
// 列出域名
result, err := provider.DomainService().List(ctx, &model.CommonPageQuery{})
domains := result.Data

// 获取域名
domain, err := provider.DomainService().Get(ctx, "example.com")

// 添加域名
newDomain := &model.Domain{
    Name:        "example.com",
    EnableHTTPS: model.EnableHTTPSOn,
}
domain, err = provider.DomainService().Add(ctx, newDomain)

// 更新域名
domain.EnableHTTPS = model.EnableHTTPSForce
domain, err = provider.DomainService().Update(ctx, domain)

// 删除域名
err = provider.DomainService().Delete(ctx, "example.com")
```

### 路由管理

**Java:**

```java
// 创建路由
Route route = new Route();
route.setName("my-route");
route.setPath("/api/v1");
route.setDomains(Arrays.asList("example.com"));

RoutePredicate predicate = new RoutePredicate();
predicate.setPathPrefix("/api/v1");
route.setPredicates(predicate);

route = provider.routeService().add(route);
```

**Go:**

```go
// 创建路由
route := &model.Route{
    Name:    "my-route",
    Path:    "/api/v1",
    Domains: []string{"example.com"},
}

route, err = provider.RouteService().Add(ctx, route)
```

### WASM 插件配置

**Java:**

```java
// 获取插件实例
WasmPluginInstance instance = provider.wasmPluginInstanceService()
    .createEmptyInstance("ai-proxy");

instance.setScope(WasmPluginInstanceScope.ROUTE);
instance.setTarget("my-route");
instance.setConfig(Map.of(
    "provider", Map.of(
        "type", "openai",
        "apiTokens", Arrays.asList("sk-xxx")
    )
));

instance = provider.wasmPluginInstanceService().addOrUpdate(instance);
```

**Go:**

```go
// 获取插件实例
instance, err := provider.WasmPluginInstanceService().CreateEmptyInstance(ctx, "ai-proxy")

instance.Scope = model.WasmPluginInstanceScopeRoute
instance.Target = "my-route"
instance.Config = map[string]interface{}{
    "provider": map[string]interface{}{
        "type":      "openai",
        "apiTokens": []string{"sk-xxx"},
    },
}

instance, err = provider.WasmPluginInstanceService().AddOrUpdate(ctx, instance)
```

### LLM 提供商配置

**Java:**

```java
LlmProvider provider = new LlmProvider();
provider.setName("my-openai");
provider.setType("openai");
provider.setApiTokens(Arrays.asList("sk-xxx"));

provider = provider.llmProviderService().addOrUpdate(provider);
```

**Go:**

```go
llmProvider := &model.LlmProvider{
    Name:      "my-openai",
    Type:      "openai",
    ApiTokens: []string{"sk-xxx"},
}

llmProvider, err = provider.LlmProviderService().AddOrUpdate(ctx, llmProvider)
```

---

## 类型映射

### 基本类型

| Java 类型 | Go 类型 |
|-----------|---------|
| `String` | `string` |
| `Integer` / `int` | `int` |
| `Long` / `long` | `int64` |
| `Boolean` / `boolean` | `bool` |
| `Double` / `double` | `float64` |
| `Object` | `interface{}` |

### 集合类型

| Java 类型 | Go 类型 |
|-----------|---------|
| `List<T>` | `[]T` (slice) |
| `Map<K,V>` | `map[K]V` |
| `Set<T>` | `map[T]struct{}` 或 `[]T` |

### 特殊类型

| Java 类型 | Go 类型 |
|-----------|---------|
| `Optional<T>` | `*T` (指针) |
| `null` | `nil` |
| `Enum` | `const` 或 `string` |

### 时间类型

| Java 类型 | Go 类型 |
|-----------|---------|
| `Date` | `time.Time` |
| `LocalDateTime` | `time.Time` |
| `Instant` | `time.Time` |

---

## 异常处理差异

### Java 异常处理

```java
try {
    Domain domain = provider.domainService().get("nonexistent.com");
} catch (NotFoundException e) {
    System.out.println("域名不存在: " + e.getMessage());
} catch (ResourceConflictException e) {
    System.out.println("资源冲突: " + e.getMessage());
} catch (BusinessException e) {
    System.out.println("业务错误: " + e.getMessage());
}
```

### Go 错误处理

```go
import "github.com/Jayj1997/higress-admin-sdk-golang/pkg/errors"

domain, err := provider.DomainService().Get(ctx, "nonexistent.com")
if err != nil {
    if errors.IsNotFoundError(err) {
        fmt.Println("域名不存在:", err.Error())
    } else if errors.IsResourceConflictError(err) {
        fmt.Println("资源冲突:", err.Error())
    } else if errors.IsBusinessError(err) {
        fmt.Println("业务错误:", err.Error())
    }
    return
}
```

### 错误类型对照

| Java 异常 | Go 错误 | 检查方法 |
|-----------|---------|----------|
| `BusinessException` | `BusinessError` | `errors.IsBusinessError(err)` |
| `NotFoundException` | `NotFoundError` | `errors.IsNotFoundError(err)` |
| `ResourceConflictException` | `ResourceConflictError` | `errors.IsResourceConflictError(err)` |
| `ValidationException` | `ValidationError` | `errors.IsValidationError(err)` |

---

## 迁移检查清单

- [ ] 更新包导入路径
- [ ] 修改配置初始化代码
- [ ] 添加 context 参数到所有服务调用
- [ ] 处理 error 返回值
- [ ] 将 null 检查改为 nil 检查
- [ ] 将 List 改为 slice
- [ ] 将 Map 改为 map
- [ ] 更新异常处理逻辑
- [ ] 测试所有功能

---

## 常见迁移问题

### Q: Go SDK 为什么需要传入 context？

A: Go 的 context 包用于控制请求的生命周期，支持超时、取消等操作。这是 Go 语言的最佳实践。

### Q: 如何处理 Java 中的 Optional？

A: Go 中使用指针表示可选值。如果值为 nil，表示不存在。

### Q: Java 中的 Builder 模式在 Go 中如何实现？

A: Go SDK 使用函数选项模式 (Functional Options Pattern)：

```go
cfg := config.NewHigressServiceConfig(
    config.WithKubeConfigPath("/path/to/kubeconfig"),
    config.WithControllerNamespace("higress-system"),
)
```

### Q: 如何处理 Java 中的异步操作？

A: Go 使用 goroutine 处理并发：

```go
go func() {
    domains, err := provider.DomainService().List(ctx, &model.CommonPageQuery{})
    // 处理结果
}()
```

---

## 相关链接

- [使用指南](usage.md)
- [API 文档](https://pkg.go.dev/github.com/Jayj1997/higress-admin-sdk-golang)
- [示例代码](../examples/)
- [Higress 官方文档](https://higress.io/zh-cn/docs/latest/overview/what-is-higress)