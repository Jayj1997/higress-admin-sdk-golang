// Example: ai-route
// This example demonstrates how to manage AI routes and LLM providers using the Higress Admin SDK.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/client"
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/config"
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model"
)

func main() {
	// Get kubeconfig path from environment variable or use default
	kubeConfigPath := os.Getenv("KUBECONFIG")
	if kubeConfigPath == "" {
		kubeConfigPath = os.Getenv("HOME") + "/.kube/config"
	}

	// Create configuration
	cfg := config.NewHigressServiceConfig(
		config.WithKubeConfigPath(kubeConfigPath),
		config.WithControllerNamespace("higress-system"),
	)

	// Create service provider
	provider, err := client.NewHigressServiceProvider(cfg)
	if err != nil {
		log.Fatalf("Failed to create provider: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("=== AI Route Management Example ===")
	fmt.Println()

	// Step 1: List existing LLM providers
	fmt.Println("Step 1: Listing existing LLM providers...")
	providers, err := provider.LlmProviderService().List(ctx)
	if err != nil {
		log.Printf("Failed to list LLM providers: %v", err)
	} else {
		fmt.Printf("Found %d LLM providers\n", len(providers))
		for i, p := range providers {
			fmt.Printf("  %d. %s (type: %s)\n", i+1, p.Name, p.Type)
		}
	}
	fmt.Println()

	// Step 2: Create an OpenAI provider
	fmt.Println("Step 2: Creating an OpenAI provider...")
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

	createdOpenai, err := provider.LlmProviderService().Add(ctx, openaiProvider)
	if err != nil {
		log.Printf("Failed to create OpenAI provider: %v", err)
	} else {
		fmt.Printf("Created LLM provider: %s\n", createdOpenai.Name)
	}
	fmt.Println()

	// Step 3: Create a Qwen (通义千问) provider
	fmt.Println("Step 3: Creating a Qwen provider...")
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

	createdQwen, err := provider.LlmProviderService().Add(ctx, qwenProvider)
	if err != nil {
		log.Printf("Failed to create Qwen provider: %v", err)
	} else {
		fmt.Printf("Created LLM provider: %s\n", createdQwen.Name)
	}
	fmt.Println()

	// Step 4: Create an Azure OpenAI provider
	fmt.Println("Step 4: Creating an Azure OpenAI provider...")
	azureProvider := &model.LlmProvider{
		Name:     "my-azure",
		Type:     "azure",
		Protocol: model.LlmProviderProtocolOpenaiV1,
		Tokens:   []string{"your-azure-api-key"},
		RawConfigs: map[string]interface{}{
			"azureBaseUrl": "https://your-resource.openai.azure.com",
		},
	}

	createdAzure, err := provider.LlmProviderService().Add(ctx, azureProvider)
	if err != nil {
		log.Printf("Failed to create Azure provider: %v", err)
	} else {
		fmt.Printf("Created LLM provider: %s\n", createdAzure.Name)
	}
	fmt.Println()

	// Step 5: Get a specific provider
	fmt.Println("Step 5: Getting a specific provider...")
	p, err := provider.LlmProviderService().Get(ctx, "my-openai")
	if err != nil {
		log.Printf("Failed to get provider: %v", err)
	} else {
		fmt.Printf("Provider name: %s\n", p.Name)
		fmt.Printf("Provider type: %s\n", p.Type)
		fmt.Printf("Provider protocol: %s\n", p.Protocol)
		if p.TokenFailoverConfig != nil {
			fmt.Printf("Token failover enabled: %v\n", p.TokenFailoverConfig.Enabled)
		}
	}
	fmt.Println()

	// Step 6: List all providers again
	fmt.Println("Step 6: Listing all providers after creation...")
	allProviders, err := provider.LlmProviderService().List(ctx)
	if err != nil {
		log.Printf("Failed to list providers: %v", err)
	} else {
		fmt.Printf("Found %d LLM providers\n", len(allProviders))
		for i, pp := range allProviders {
			fmt.Printf("  %d. %s (type: %s)\n", i+1, pp.Name, pp.Type)
		}
	}
	fmt.Println()

	// Step 7: Delete the providers
	fmt.Println("Step 7: Cleaning up - deleting providers...")
	err = provider.LlmProviderService().Delete(ctx, "my-openai")
	if err != nil {
		log.Printf("Failed to delete OpenAI provider: %v", err)
	} else {
		fmt.Println("Deleted OpenAI provider")
	}

	err = provider.LlmProviderService().Delete(ctx, "my-qwen")
	if err != nil {
		log.Printf("Failed to delete Qwen provider: %v", err)
	} else {
		fmt.Println("Deleted Qwen provider")
	}

	err = provider.LlmProviderService().Delete(ctx, "my-azure")
	if err != nil {
		log.Printf("Failed to delete Azure provider: %v", err)
	} else {
		fmt.Println("Deleted Azure provider")
	}
	fmt.Println()

	fmt.Println("=== Example Complete ===")
}
