// Package ai provides AI-related services for Higress Admin SDK.
package ai

import (
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model"
)

// 初始化所有内置处理器
func init() {
	// 注册OpenAI处理器
	RegisterHandler(NewOpenaiLlmProviderHandler())

	// 注册通义千问处理器
	RegisterHandler(NewQwenLlmProviderHandler())

	// 注册Azure处理器
	RegisterHandler(NewAzureLlmProviderHandler())

	// 注册Ollama处理器
	RegisterHandler(NewOllamaLlmProviderHandler())

	// 注册Moonshot处理器
	RegisterHandler(NewDefaultLlmProviderHandler(
		model.LlmProviderTypeMoonshot,
		"api.moonshot.cn",
		443,
		"https",
	))

	// 注册AI360处理器
	RegisterHandler(NewDefaultLlmProviderHandler(
		model.LlmProviderTypeAi360,
		"api.360.cn",
		443,
		"https",
	))

	// 注册GitHub处理器
	RegisterHandler(NewDefaultLlmProviderHandler(
		model.LlmProviderTypeGithub,
		"models.inference.ai.azure.com",
		443,
		"https",
	))

	// 注册Groq处理器
	RegisterHandler(NewDefaultLlmProviderHandler(
		model.LlmProviderTypeGroq,
		"api.groq.com",
		443,
		"https",
	))

	// 注册百川处理器
	RegisterHandler(NewDefaultLlmProviderHandler(
		model.LlmProviderTypeBaichuan,
		"api.baichuan-ai.com",
		443,
		"https",
	))

	// 注册零一万物处理器
	RegisterHandler(NewDefaultLlmProviderHandler(
		model.LlmProviderTypeYi,
		"api.lingyiwanwu.com",
		443,
		"https",
	))

	// 注册DeepSeek处理器
	RegisterHandler(NewDefaultLlmProviderHandler(
		model.LlmProviderTypeDeepSeek,
		"api.deepseek.com",
		443,
		"https",
	))

	// 注册智谱AI处理器
	RegisterHandler(NewDefaultLlmProviderHandler(
		model.LlmProviderTypeZhipuai,
		"open.bigmodel.cn",
		443,
		"https",
	))

	// 注册Claude处理器
	RegisterHandler(NewDefaultLlmProviderHandler(
		model.LlmProviderTypeClaude,
		"api.anthropic.com",
		443,
		"https",
	))

	// 注册百度处理器
	RegisterHandler(NewDefaultLlmProviderHandler(
		model.LlmProviderTypeBaidu,
		"qianfan.baidubce.com",
		443,
		"https",
	))

	// 注册阶跃星辰处理器
	RegisterHandler(NewDefaultLlmProviderHandler(
		model.LlmProviderTypeStepfun,
		"api.stepfun.com",
		443,
		"https",
	))

	// 注册MiniMax处理器
	RegisterHandler(NewDefaultLlmProviderHandler(
		model.LlmProviderTypeMinimax,
		"api.minimax.chat",
		443,
		"https",
	))

	// 注册Gemini处理器
	RegisterHandler(NewDefaultLlmProviderHandler(
		model.LlmProviderTypeGemini,
		"generativelanguage.googleapis.com",
		443,
		"https",
	))

	// 注册Mistral处理器
	RegisterHandler(NewDefaultLlmProviderHandler(
		model.LlmProviderTypeMistral,
		"api.mistral.ai",
		443,
		"https",
	))

	// 注册Cohere处理器
	RegisterHandler(NewDefaultLlmProviderHandler(
		model.LlmProviderTypeCohere,
		"api.cohere.com",
		443,
		"https",
	))

	// 注册豆包处理器
	RegisterHandler(NewDefaultLlmProviderHandler(
		model.LlmProviderTypeDoubao,
		"ark.cn-beijing.volces.com",
		443,
		"https",
	))

	// 注册Coze处理器
	RegisterHandler(NewDefaultLlmProviderHandler(
		model.LlmProviderTypeCoze,
		"api.coze.cn",
		443,
		"https",
	))

	// 注册OpenRouter处理器
	RegisterHandler(NewDefaultLlmProviderHandler(
		model.LlmProviderTypeOpenrouter,
		"openrouter.ai",
		443,
		"https",
	))

	// 注册Grok处理器
	RegisterHandler(NewDefaultLlmProviderHandler(
		model.LlmProviderTypeGrok,
		"api.x.ai",
		443,
		"https",
	))

	// 注册Bedrock处理器
	RegisterHandler(NewBedrockLlmProviderHandler())

	// 注册Vertex处理器
	RegisterHandler(NewVertexLlmProviderHandler())
}
