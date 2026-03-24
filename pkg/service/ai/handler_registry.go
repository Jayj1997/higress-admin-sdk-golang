// Package ai provides AI-related services for Higress Admin SDK.
package ai

import (
	"sync"
)

var (
	handlerRegistry = make(map[string]LlmProviderHandler)
	registryMutex   sync.RWMutex
)

// RegisterHandler 注册处理器
func RegisterHandler(handler LlmProviderHandler) {
	registryMutex.Lock()
	defer registryMutex.Unlock()
	handlerRegistry[handler.GetType()] = handler
}

// GetHandler 获取处理器
func GetHandler(providerType string) LlmProviderHandler {
	registryMutex.RLock()
	defer registryMutex.RUnlock()
	return handlerRegistry[providerType]
}

// GetAllHandlers 获取所有处理器
func GetAllHandlers() map[string]LlmProviderHandler {
	registryMutex.RLock()
	defer registryMutex.RUnlock()
	result := make(map[string]LlmProviderHandler, len(handlerRegistry))
	for k, v := range handlerRegistry {
		result[k] = v
	}
	return result
}

// HasHandler 检查处理器是否存在
func HasHandler(providerType string) bool {
	registryMutex.RLock()
	defer registryMutex.RUnlock()
	_, exists := handlerRegistry[providerType]
	return exists
}
