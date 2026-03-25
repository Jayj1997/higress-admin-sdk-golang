# Route Management Example

本示例演示如何使用 Higress Admin SDK 管理网关路由。

## 功能演示

1. **列出路由** - 分页查询现有路由
2. **创建路由** - 创建带有路径匹配、CORS、请求头控制的路由
3. **获取路由** - 根据名称获取路由详情
4. **更新路由** - 修改路由配置（添加重写规则等）
5. **过滤查询** - 按域名过滤路由列表
6. **删除路由** - 删除指定路由

## 运行方式

```bash
# 设置 kubeconfig 路径（可选）
export KUBECONFIG=~/.kube/config

# 运行示例
go run main.go
```

## 代码说明

### 创建路由

```go
newRoute := &model.Route{
    Name:    "example-route",
    Domains: []string{"example.com"},
    Path: &route.RoutePredicate{
        MatchType: route.MatchTypePrefix,
        Path:      "/api/example",
    },
    Services: []*route.UpstreamService{
        {
            Name:      "example-service",
            Namespace: "default",
            Port:      8080,
        },
    },
}
```

### 配置 CORS

```go
newRoute.CORS = &route.CorsConfig{
    AllowOrigins:     []string{"*"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Content-Type", "Authorization"},
    AllowCredentials: &allowCredentials,
    MaxAge:           &maxAge,
}
```

### 配置请求头控制

```go
newRoute.HeaderControl = &route.HeaderControlConfig{
    RequestAddHeaders: map[string]string{
        "X-Custom-Header": "example-value",
    },
}
```

### 配置路径重写

```go
route.Rewrite = &route.RewriteConfig{
    Path: "/v2/api/example",
}
```

## 注意事项

- 运行示例需要连接到 Kubernetes 集群
- 示例会在最后删除创建的路由
- 请确保 Higress 控制器已部署在指定命名空间