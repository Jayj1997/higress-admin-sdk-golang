# AI Route Management Example

本示例演示如何使用 Higress Admin SDK 管理 AI 路由和 LLM 提供商。

## 功能演示

1. **列出 LLM 提供商** - 查看所有已配置的 LLM 提供商
2. **创建 OpenAI 提供商** - 配置 OpenAI API 连接
3. **创建通义千问提供商** - 配置阿里云通义千问 API 连接
4. **创建 Azure OpenAI 提供商** - 配置 Azure OpenAI API 连接
5. **获取提供商详情** - 根据名称获取提供商配置
6. **删除提供商** - 清理创建的提供商

## 运行方式

```bash
# 设置 kubeconfig 路径（可选）
export KUBECONFIG=~/.kube/config

# 运行示例
go run main.go
```

## 支持的 LLM 提供商

SDK 支持以下 LLM 提供商：

| 提供商 | 类型标识 | 说明 |
|--------|----------|------|
| OpenAI | `openai` | OpenAI 官方 API |
| Azure OpenAI | `azure` | Azure OpenAI 服务 |
| 通义千问 | `qwen` | 阿里云通义千问 |
| 文心一言 | `ernie` | 百度文心一言 |
| 智谱 AI | `glm` | 智谱 ChatGLM |
| Ollama | `ollama` | 本地 Ollama 服务 |
| AWS Bedrock | `bedrock` | AWS Bedrock 服务 |
| GCP Vertex AI | `vertex` | Google Vertex AI |
| Moonshot | `moonshot` | Moonshot AI |
| DeepSeek | `deepseek` | DeepSeek AI |
| Yi | `yi` | 零一万物 |
| Baichuan | `baichuan` | 百川智能 |
| Minimax | `minimax` | Minimax |
| Claude | `claude` | Anthropic Claude |

## 代码说明

### 创建 OpenAI 提供商

```go
openaiProvider := &model.LlmProvider{
    Name:     "my-openai",
    Type:     "openai",
    Protocol: model.LlmProviderProtocolOpenaiV1,
    Tokens:   []string{"sk-your-openai-api-key"},
    TokenFailoverConfig: &model.TokenFailoverConfig{
        Enabled:             true,
        FailureThreshold:    3,
        SuccessThreshold:    1,
        HealthCheckInterval: 300,
        HealthCheckTimeout:  30,
    },
}
```

### 创建通义千问提供商

```go
qwenProvider := &model.LlmProvider{
    Name:     "my-qwen",
    Type:     "qwen",
    Protocol: model.LlmProviderProtocolOpenaiV1,
    Tokens:   []string{"your-qwen-api-key"},
    RawConfigs: map[string]interface{}{
        "modelMapping": map[string]interface{}{
            "qwen-turbo": "qwen-turbo",
            "qwen-plus":  "qwen-plus",
        },
    },
}
```

### 创建 Azure OpenAI 提供商

```go
azureProvider := &model.LlmProvider{
    Name:     "my-azure",
    Type:     "azure",
    Protocol: model.LlmProviderProtocolOpenaiV1,
    Tokens:   []string{"your-azure-api-key"},
    RawConfigs: map[string]interface{}{
        "azureBaseUrl": "https://your-resource.openai.azure.com",
    },
}
```

## 注意事项

- 运行示例需要连接到 Kubernetes 集群
- 示例会在最后删除创建的提供商
- 请确保 Higress 控制器已部署在指定命名空间
- Token 故障转移配置是可选的，用于自动切换失效的 API Key