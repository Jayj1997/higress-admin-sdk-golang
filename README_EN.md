# Higress Admin SDK for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/Jayj1997/higress-admin-sdk-golang/v2.svg)](https://pkg.go.dev/github.com/Jayj1997/higress-admin-sdk-golang/v2)
[![Version](https://img.shields.io/badge/version-v2.2.0-blue.svg)](https://github.com/Jayj1997/higress-admin-sdk-golang/v2/releases)

English | [简体中文](README.md)

A Go SDK for managing Higress gateway configurations, including domains, routes, services, certificates, and WASM plugins.

> **Version Note**: This project is ported from the official higress-admin-sdk, version **v2.2.0** is aligned with [higress-console 2.2.0](https://github.com/higress-group/higress-console).

## Features

- **Domain Management**: Create, update, delete, and list gateway domains
- **Route Management**: Configure routing rules with path matching, header control, CORS, rewrite, and more
- **Service Management**: List and manage backend services
- **Service Source Management**: Configure service discovery sources (Nacos, DNS, static, etc.)
- **TLS Certificate Management**: Manage SSL/TLS certificates for HTTPS
- **WASM Plugin Management**: Configure and manage WASM plugins (40+ built-in plugins)
- **AI Route Management**: Configure AI/LLM routing rules (supports 20+ LLM providers)
- **Consumer Management**: Manage API consumers and credentials
- **MCP Server Management**: Configure MCP server instances

## Installation

```bash
go get github.com/Jayj1997/higress-admin-sdk-golang/v2
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    sdk "github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/client"
    "github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/config"
    "github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model"
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
    domains, err := provider.DomainService().List(context.Background(), &model.CommonPageQuery{})
    if err != nil {
        log.Fatalf("Failed to list domains: %v", err)
    }

    fmt.Printf("Found %d domains\n", domains.Total)
    for _, domain := range domains.Data {
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

## Documentation

- **[Usage Guide](docs/usage.md)** - Detailed configuration and service usage
- **[Migration Guide](docs/migration.md)** - Migration from Java SDK to Go SDK
- **[API Documentation](https://pkg.go.dev/github.com/Jayj1997/higress-admin-sdk-golang/v2)** - GoDoc generated API reference
- **[Changelog](CHANGELOG.md)** - Version history

## Examples

| Example | Description |
|---------|-------------|
| [hello-sdk](examples/hello-sdk/) | Basic usage example |
| [route-management](examples/route-management/) | Route management (CRUD, CORS, header control) |
| [ai-route](examples/ai-route/) | AI route management (LLM provider management) |
| [wasm-plugin](examples/wasm-plugin/) | WASM plugin configuration |

## API Reference

### Domain Service

```go
// List all domains
domains, err := provider.DomainService().List(ctx, &model.CommonPageQuery{})

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
    Domains: []string{"example.com"},
    Path: &route.RoutePredicate{
        MatchType: route.MatchTypePrefix,
        Path:      "/api/v1",
    },
    Services: []*route.UpstreamService{
        {Name: "my-service", Namespace: "default", Port: 8080},
    },
}
route, err := provider.RouteService().Add(ctx, newRoute)
```

### WASM Plugin Service

```go
// List built-in plugins
builtIn := true
plugins, err := provider.WasmPluginService().List(ctx, &model.WasmPluginPageQuery{
    BuiltIn: &builtIn,
})

// Get plugin configuration
config, err := provider.WasmPluginService().GetConfig(ctx, "ai-proxy", "en-US")

// Create plugin instance
instance, err := provider.WasmPluginInstanceService().CreateEmptyInstance(ctx, "ai-proxy")
instance.Scope = model.WasmPluginInstanceScopeRoute
instance.Target = "my-route"
instance.Configurations = map[string]interface{}{
    "provider": map[string]interface{}{
        "type":      "openai",
        "apiTokens": []string{"sk-xxx"},
    },
}
instance, err = provider.WasmPluginInstanceService().AddOrUpdate(ctx, instance)
```

## Development

### Prerequisites

- Go 1.21+
- Make

### Running Tests

```bash
# Run unit tests (no K8s cluster required)
make test

# Run unit tests (alias for make test)
make test-unit

# Run integration tests (requires K8s cluster)
make test-integration

# Run all tests including integration tests (requires K8s cluster)
make test-all

# Run tests with coverage
make test-coverage

# Run linting
make lint
```

#### Test Types

| Command | Description | K8s Required |
|---------|-------------|--------------|
| `make test` | Unit tests only | No |
| `make test-unit` | Same as `make test` | No |
| `make test-integration` | Integration tests only | Yes |
| `make test-all` | All tests (unit + integration) | Yes |

### Project Structure

```
higress-admin-sdk-golang/
├── CHANGELOG.md           # Changelog
├── docs/
│   ├── usage.md           # Usage guide
│   └── migration.md       # Migration guide
├── examples/
│   ├── hello-sdk/         # Basic example
│   ├── route-management/  # Route management example
│   ├── ai-route/          # AI route example
│   └── wasm-plugin/       # WASM plugin example
├── pkg/
│   ├── client/            # Client implementation
│   ├── config/            # Configuration
│   ├── model/             # Data models
│   ├── service/           # Service implementations
│   ├── constant/          # Constants
│   └── errors/            # Error definitions
├── internal/
│   ├── kubernetes/        # Kubernetes client and CRD
│   └── resources/         # Built-in resources (WASM plugin configs)
└── test/                  # Integration tests
```

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for details.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Related Projects

- [Higress](https://github.com/alibaba/higress) - Cloud-native API Gateway
- [Higress Console](https://github.com/higress-group/higress-console) - Console for Higress