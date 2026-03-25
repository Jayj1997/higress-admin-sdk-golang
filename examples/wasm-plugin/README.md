# WASM Plugin Management Example

本示例演示如何使用 Higress Admin SDK 管理 WASM 插件。

## 功能演示

1. **列出内置插件** - 查看所有内置的 WASM 插件
2. **获取插件详情** - 获取特定插件的详细信息
3. **获取配置 Schema** - 获取插件的配置定义
4. **获取 README 文档** - 获取插件的使用文档
5. **创建插件实例** - 为路由配置插件实例
6. **列出插件实例** - 查看特定插件的所有实例
7. **查询插件实例** - 查询特定范围的插件配置
8. **删除插件实例** - 清理插件配置

## 运行方式

```bash
# 设置 kubeconfig 路径（可选）
export KUBECONFIG=~/.kube/config

# 运行示例
go run main.go
```

## 插件作用域

WASM 插件实例可以在不同级别配置：

| 作用域 | 说明 |
|--------|------|
| `global` | 全局配置，应用于所有请求 |
| `domain` | 域名级别，应用于特定域名的请求 |
| `route` | 路由级别，应用于特定路由的请求 |
| `service` | 服务级别，应用于特定服务的请求 |

## 内置插件列表

SDK 内置了 40+ 个插件，包括：

### AI 相关
- `ai-proxy` - AI 代理插件
- `ai-statistics` - AI 统计插件
- `ai-token-ratelimit` - AI Token 限流
- `ai-cache` - AI 缓存
- `ai-rag` - RAG 支持

### 安全相关
- `ai-security-guard` - AI 安全防护
- `ai-data-masking` - 数据脱敏

### 其他
- `request-validation` - 请求验证
- `key-rate-limit` - 限流
- `basic-auth` - 基础认证
- `jwt-auth` - JWT 认证
- `ip-restriction` - IP 限制

## 代码说明

### 列出内置插件

```go
builtInTrue := true
plugins, err := provider.WasmPluginService().List(ctx, &model.WasmPluginPageQuery{
    BuiltIn: &builtInTrue,
})
```

### 创建路由级插件实例

```go
instance, err := provider.WasmPluginInstanceService().CreateEmptyInstance(ctx, "ai-proxy")
instance.Scope = model.WasmPluginInstanceScopeRoute
instance.Target = "my-route"
instance.Configurations = map[string]interface{}{
    "provider": map[string]interface{}{
        "type":      "openai",
        "apiTokens": []string{"sk-your-api-key"},
    },
}
instance, err = provider.WasmPluginInstanceService().AddOrUpdate(ctx, instance)
```

### 查询插件实例

```go
instance, err := provider.WasmPluginInstanceService().Query(ctx,
    model.WasmPluginInstanceScopeRoute,
    "my-route",
    "ai-proxy",
    nil,
)
```

## 注意事项

- 运行示例需要连接到 Kubernetes 集群
- 示例会在最后删除创建的插件实例
- 请确保 Higress 控制器已部署在指定命名空间
- 插件配置格式请参考各插件的 Schema 定义