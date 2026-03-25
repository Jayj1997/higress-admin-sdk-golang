// Example: route-management
// This example demonstrates how to manage routes using the Higress Admin SDK.
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
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model/route"
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

	fmt.Println("=== Route Management Example ===")
	fmt.Println()

	// Step 1: List existing routes
	fmt.Println("Step 1: Listing existing routes...")
	routes, err := provider.RouteService().List(ctx, &model.RoutePageQuery{
		CommonPageQuery: model.CommonPageQuery{
			PageNum:  1,
			PageSize: 10,
		},
	})
	if err != nil {
		log.Printf("Failed to list routes: %v", err)
	} else {
		fmt.Printf("Found %d routes\n", routes.Total)
		for i, r := range routes.Data {
			fmt.Printf("  %d. %s\n", i+1, r.Name)
		}
	}
	fmt.Println()

	// Step 2: Create a new route
	fmt.Println("Step 2: Creating a new route...")
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

	// Add header control configuration
	newRoute.HeaderControl = &route.HeaderControlConfig{
		RequestAddHeaders: map[string]string{
			"X-Custom-Header": "example-value",
		},
	}

	// Add CORS configuration
	allowCredentials := true
	maxAge := 3600
	newRoute.CORS = &route.CorsConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: &allowCredentials,
		MaxAge:           &maxAge,
	}

	createdRoute, err := provider.RouteService().Add(ctx, newRoute)
	if err != nil {
		log.Printf("Failed to create route: %v", err)
	} else {
		fmt.Printf("Created route: %s\n", createdRoute.Name)
	}
	fmt.Println()

	// Step 3: Get the route
	fmt.Println("Step 3: Getting the created route...")
	r, err := provider.RouteService().Get(ctx, "example-route")
	if err != nil {
		log.Printf("Failed to get route: %v", err)
	} else {
		fmt.Printf("Route name: %s\n", r.Name)
		if r.Path != nil {
			fmt.Printf("Route path: %s (match: %s)\n", r.Path.Path, r.Path.MatchType)
		}
		fmt.Printf("Route domains: %v\n", r.Domains)
		if r.CORS != nil {
			fmt.Printf("CORS allow origins: %v\n", r.CORS.AllowOrigins)
		}
	}
	fmt.Println()

	// Step 4: Update the route
	fmt.Println("Step 4: Updating the route...")
	if r != nil {
		// Add path rewrite
		r.Rewrite = &route.RewriteConfig{
			Path: "/v2/api/example",
		}

		// Update header control
		if r.HeaderControl == nil {
			r.HeaderControl = &route.HeaderControlConfig{}
		}
		if r.HeaderControl.RequestAddHeaders == nil {
			r.HeaderControl.RequestAddHeaders = make(map[string]string)
		}
		r.HeaderControl.RequestAddHeaders["X-Updated"] = "true"

		updatedRoute, err := provider.RouteService().Update(ctx, r)
		if err != nil {
			log.Printf("Failed to update route: %v", err)
		} else {
			fmt.Printf("Updated route: %s\n", updatedRoute.Name)
			if updatedRoute.Rewrite != nil {
				fmt.Printf("Rewrite path: %s\n", updatedRoute.Rewrite.Path)
			}
		}
	}
	fmt.Println()

	// Step 5: List routes with domain filter
	fmt.Println("Step 5: Listing routes filtered by domain...")
	filteredRoutes, err := provider.RouteService().List(ctx, &model.RoutePageQuery{
		CommonPageQuery: model.CommonPageQuery{
			PageNum:  1,
			PageSize: 10,
		},
		DomainName: "example.com",
	})
	if err != nil {
		log.Printf("Failed to list routes: %v", err)
	} else {
		fmt.Printf("Found %d routes for domain 'example.com'\n", filteredRoutes.Total)
		for i, rr := range filteredRoutes.Data {
			fmt.Printf("  %d. %s\n", i+1, rr.Name)
		}
	}
	fmt.Println()

	// Step 6: Delete the route
	fmt.Println("Step 6: Deleting the route...")
	err = provider.RouteService().Delete(ctx, "example-route")
	if err != nil {
		log.Printf("Failed to delete route: %v", err)
	} else {
		fmt.Println("Route deleted successfully")
	}
	fmt.Println()

	fmt.Println("=== Example Complete ===")
}
