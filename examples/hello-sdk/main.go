// Example: hello-sdk
// This example demonstrates how to use the Higress Admin SDK to interact with Higress resources.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

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

	fmt.Println("Higress Admin SDK initialized successfully!")
	fmt.Println()

	// List all routes
	fmt.Println("=== Listing Routes ===")
	routes, err := provider.RouteService().List(context.Background(), &model.RoutePageQuery{})
	if err != nil {
		log.Printf("Failed to list routes: %v", err)
	} else {
		fmt.Printf("Found %d routes\n", routes.Total)
		for i, route := range routes.Data {
			fmt.Printf("  %d. %s\n", i+1, route.Name)
		}
	}
	fmt.Println()

	// List all domains
	fmt.Println("=== Listing Domains ===")
	domains, err := provider.DomainService().List(context.Background(), &model.CommonPageQuery{})
	if err != nil {
		log.Printf("Failed to list domains: %v", err)
	} else {
		fmt.Printf("Found %d domains\n", domains.Total)
		for i, domain := range domains.Data {
			fmt.Printf("  %d. %s\n", i+1, domain.Name)
		}
	}
	fmt.Println()

	// List all services
	fmt.Println("=== Listing Services ===")
	services, err := provider.ServiceService().List(context.Background(), &model.CommonPageQuery{})
	if err != nil {
		log.Printf("Failed to list services: %v", err)
	} else {
		fmt.Printf("Found %d services\n", services.Total)
		for i, svc := range services.Data {
			fmt.Printf("  %d. %s/%s\n", i+1, svc.Namespace, svc.Name)
		}
	}
	fmt.Println()

	// List all service sources
	fmt.Println("=== Listing Service Sources ===")
	sources, err := provider.ServiceSourceService().List(context.Background(), &model.CommonPageQuery{})
	if err != nil {
		log.Printf("Failed to list service sources: %v", err)
	} else {
		fmt.Printf("Found %d service sources\n", sources.Total)
		for i, source := range sources.Data {
			fmt.Printf("  %d. %s (type: %s)\n", i+1, source.Name, source.Type)
		}
	}
	fmt.Println()

	// List all TLS certificates
	fmt.Println("=== Listing TLS Certificates ===")
	certs, err := provider.TlsCertificateService().List(context.Background(), &model.CommonPageQuery{})
	if err != nil {
		log.Printf("Failed to list TLS certificates: %v", err)
	} else {
		fmt.Printf("Found %d TLS certificates\n", certs.Total)
		for i, cert := range certs.Data {
			fmt.Printf("  %d. %s\n", i+1, cert.Name)
		}
	}

	fmt.Println()
	fmt.Println("Done!")
}
