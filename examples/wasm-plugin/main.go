// Example: wasm-plugin
// This example demonstrates how to manage WASM plugins using the Higress Admin SDK.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/client"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/config"
	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model"
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

	fmt.Println("=== WASM Plugin Management Example ===")
	fmt.Println()

	// Step 1: List built-in plugins
	fmt.Println("Step 1: Listing built-in WASM plugins...")
	builtInTrue := true
	plugins, err := provider.WasmPluginService().List(ctx, &model.WasmPluginPageQuery{
		BuiltIn: &builtInTrue,
	})
	if err != nil {
		log.Printf("Failed to list plugins: %v", err)
	} else {
		fmt.Printf("Found %d built-in plugins\n", plugins.Total)
		// Show first 10 plugins
		count := 0
		for _, p := range plugins.Data {
			if count >= 10 {
				break
			}
			fmt.Printf("  %d. %s\n", count+1, p.Name)
			count++
		}
		if plugins.Total > 10 {
			fmt.Printf("  ... and %d more\n", plugins.Total-10)
		}
	}
	fmt.Println()

	// Step 2: Get plugin details
	fmt.Println("Step 2: Getting ai-proxy plugin details...")
	plugin, err := provider.WasmPluginService().Get(ctx, "ai-proxy", "zh-CN")
	if err != nil {
		log.Printf("Failed to get plugin: %v", err)
	} else {
		fmt.Printf("Plugin name: %s\n", plugin.Name)
		fmt.Printf("Plugin version: %s\n", plugin.PluginVersion)
		if plugin.Category != "" {
			fmt.Printf("Plugin category: %s\n", plugin.Category)
		}
	}
	fmt.Println()

	// Step 3: Get plugin config schema
	fmt.Println("Step 3: Getting ai-proxy plugin config schema...")
	pluginConfig, err := provider.WasmPluginService().GetConfig(ctx, "ai-proxy", "zh-CN")
	if err != nil {
		log.Printf("Failed to get plugin config: %v", err)
	} else {
		fmt.Printf("Config schema available: %v\n", pluginConfig.Schema != nil)
	}
	fmt.Println()

	// Step 4: Get plugin README
	fmt.Println("Step 4: Getting ai-proxy plugin README...")
	readme, err := provider.WasmPluginService().GetReadme(ctx, "ai-proxy", "zh-CN")
	if err != nil {
		log.Printf("Failed to get plugin readme: %v", err)
	} else {
		// Show first 200 characters of README
		if len(readme) > 200 {
			fmt.Printf("README preview: %s...\n", readme[:200])
		} else {
			fmt.Printf("README: %s\n", readme)
		}
	}
	fmt.Println()

	// Step 5: Create a plugin instance for a route
	fmt.Println("Step 5: Creating a plugin instance for a route...")
	instance, err := provider.WasmPluginInstanceService().CreateEmptyInstance(ctx, "ai-proxy")
	if err != nil {
		log.Printf("Failed to create empty instance: %v", err)
	} else {
		instance.Scope = model.WasmPluginInstanceScopeRoute
		instance.Target = "example-route"
		instance.Configurations = map[string]interface{}{
			"provider": map[string]interface{}{
				"type":      "openai",
				"apiTokens": []string{"sk-your-api-key"},
			},
		}

		createdInstance, err := provider.WasmPluginInstanceService().AddOrUpdate(ctx, instance)
		if err != nil {
			log.Printf("Failed to create instance: %v", err)
		} else {
			fmt.Printf("Created plugin instance for route: %s\n", createdInstance.Target)
		}
	}
	fmt.Println()

	// Step 6: List plugin instances
	fmt.Println("Step 6: Listing plugin instances for ai-proxy...")
	instances, err := provider.WasmPluginInstanceService().ListByPlugin(ctx, "ai-proxy", nil)
	if err != nil {
		log.Printf("Failed to list instances: %v", err)
	} else {
		fmt.Printf("Found %d instances\n", len(instances))
		for i, inst := range instances {
			fmt.Printf("  %d. Scope: %s, Target: %s\n", i+1, inst.Scope, inst.Target)
		}
	}
	fmt.Println()

	// Step 7: Query specific instance
	fmt.Println("Step 7: Querying specific plugin instance...")
	queryInstance, err := provider.WasmPluginInstanceService().Query(ctx,
		model.WasmPluginInstanceScopeRoute,
		"example-route",
		"ai-proxy",
		nil,
	)
	if err != nil {
		log.Printf("Failed to query instance: %v", err)
	} else {
		fmt.Printf("Found instance: scope=%s, target=%s\n", queryInstance.Scope, queryInstance.Target)
	}
	fmt.Println()

	// Step 8: Delete the plugin instance
	fmt.Println("Step 8: Deleting the plugin instance...")
	err = provider.WasmPluginInstanceService().Delete(ctx,
		model.WasmPluginInstanceScopeRoute,
		"example-route",
		"ai-proxy",
		nil,
	)
	if err != nil {
		log.Printf("Failed to delete instance: %v", err)
	} else {
		fmt.Println("Plugin instance deleted successfully")
	}
	fmt.Println()

	fmt.Println("=== Example Complete ===")
}
