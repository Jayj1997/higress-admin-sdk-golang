# Higress Admin SDK for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/Jayj1997/higress-admin-sdk-golang.svg)](https://pkg.go.dev/github.com/Jayj1997/higress-admin-sdk-golang)

English | [简体中文](README_CN.md)

A Go SDK for managing Higress gateway configurations, including domains, routes, services, certificates, and WASM plugins.

> **Note**: This project is ported from the official higress-admin-sdk by higress-group. For more details, see [higress-console](https://github.com/higress-group/higress-console) and [higress](https://github.com/alibaba/higress).

## Features

- **Domain Management**: Create, update, delete, and list gateway domains
- **Route Management**: Configure routing rules with path matching, header control, and more
- **Service Management**: List and manage backend services
- **Service Source Management**: Configure service discovery sources (Nacos, DNS, static, etc.)
- **TLS Certificate Management**: Manage SSL/TLS certificates for HTTPS
- **WASM Plugin Management**: Configure and manage WASM plugins
- **AI Route Management**: Configure AI/LLM routing rules
- **Consumer Management**: Manage API consumers and credentials
- **MCP Server Management**: Configure MCP server instances

## Installation

```bash
go get github.com/Jayj1997/higress-admin-sdk-golang
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    sdk "github.com/Jayj1997/higress-admin-sdk-golang"
    "github.com/Jayj1997/higress-admin-sdk-golang/pkg/config"
)

func main() {
    // Create configuration
    cfg := config.NewHigressServiceConfig(
        config.WithKubeConfigPath("~/.kube/config"),
        config.WithControllerNamespace("higress-system"),
    )

    // Create service provider
    provider, err := sdk.NewHigressServiceProvider(cfg)
    if err != nil {
        log.Fatalf("Failed to create provider: %v", err)
    }

    // List all domains
    domains, err := provider.DomainService().List(context.Background())
    if err != nil {
        log.Fatalf("Failed to list domains: %v", err)
    }

    fmt.Printf("Found %d domains\n", len(domains))
    for _, domain := range domains {
        fmt.Printf("- %s\n", domain.Name)
    }
}
```

## Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `WithKubeConfigPath(path)` | Path to kubeconfig file | - |
| `WithKubeConfigContent(content)` | Kubeconfig content as string | - |
| `WithControllerNamespace(ns)` | Higress controller namespace | `higress-system` |
| `WithControllerServiceHost(host)` | Controller service host | - |
| `WithControllerServicePort(port)` | Controller service port | `15014` |

## API Reference

### Domain Service

```go
// List all domains
domains, err := provider.DomainService().List(ctx)

// Get a specific domain
domain, err := provider.DomainService().Get(ctx, "example.com")

// Add a new domain
newDomain := &model.Domain{
    Name:        "example.com",
    EnableHTTPS: model.EnableHTTPSOn,
}
domain, err := provider.DomainService().Add(ctx, newDomain)

// Update a domain
domain.EnableHTTPS = model.EnableHTTPSForce
domain, err := provider.DomainService().Update(ctx, domain)

// Delete a domain
err := provider.DomainService().Delete(ctx, "example.com")
```

### Route Service

```go
// List routes with pagination
routes, err := provider.RouteService().List(ctx, &model.RoutePageQuery{
    DomainName: "example.com",
})

// Get a specific route
route, err := provider.RouteService().Get(ctx, "my-route")

// Add a new route
newRoute := &model.Route{
    Name:    "my-route",
    Path:    "/api/v1",
    Domains: []string{"example.com"},
}
route, err := provider.RouteService().Add(ctx, newRoute)
```

## Documentation

- [API Documentation](https://pkg.go.dev/github.com/Jayj1997/higress-admin-sdk-golang)
- [Migration Guide from Java SDK](docs/migration.md)
- [Examples](examples/)

## Development

### Prerequisites

- Go 1.21+
- Make

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run linting
make lint
```

### Project Structure

```
higress-admin-sdk-golang/
├── api/v1/              # API definitions
├── pkg/
│   ├── client/          # Client implementation
│   ├── config/          # Configuration
│   ├── model/           # Data models
│   ├── service/         # Service implementations
│   ├── constant/        # Constants
│   └── errors/          # Error definitions
├── internal/
│   ├── kubernetes/      # Kubernetes client and CRD
│   └── util/            # Internal utilities
├── examples/            # Example code
└── test/                # Integration tests
```

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for details.

## Related Projects

- [Higress](https://github.com/alibaba/higress) - Cloud-native API Gateway
- [Higress Console](https://github.com/higress-group/higress-console) - console for Higress
