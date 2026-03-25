# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.2.0] - 2026-03-25

### Added

This is the initial release of the Go SDK, ported from higress-admin-sdk Java (higress&higress-console 2.2.0).

#### Core Services

- **Domain Management** - DomainService for managing gateway domains
  - List, Get, Add, Update, Delete operations
  - HTTPS configuration support

- **Route Management** - RouteService for configuring routing rules
  - Path matching (exact, prefix, regex)
  - Header control configuration
  - CORS configuration
  - URL rewrite configuration
  - Redirect configuration
  - Proxy next upstream configuration

- **Service Management** - ServiceService for listing backend services
  - Support for Kubernetes services
  - Support for registered services from Nacos, DNS, etc.

- **Service Source Management** - ServiceSourceService for configuring service discovery
  - Nacos 2.x support
  - DNS service discovery
  - Static service configuration

- **TLS Certificate Management** - TlsCertificateService for managing SSL/TLS certificates
  - Certificate CRUD operations
  - Domain binding support

- **Proxy Server Management** - ProxyServerService for managing proxy servers

#### WASM Plugin System

- **WasmPluginService** - Manage WASM plugins
  - 40+ built-in plugins included
  - Custom plugin support
  - Plugin configuration schema
  - i18n support (zh-CN, en-US)

- **WasmPluginInstanceService** - Configure plugin instances
  - Global scope configuration
  - Domain scope configuration
  - Route scope configuration
  - Service scope configuration

#### AI/LLM Features

- **LlmProviderService** - Manage LLM providers
  - Support for 20+ LLM providers:
    - OpenAI
    - Azure OpenAI
    - Qwen (通义千问)
    - ERNIE (文心一言)
    - GLM (智谱AI)
    - Ollama
    - AWS Bedrock
    - GCP Vertex AI
    - Moonshot
    - DeepSeek
    - Yi (零一万物)
    - Baichuan (百川智能)
    - Minimax
    - Claude
    - And more...
  - Token failover configuration
  - Model mapping support

- **AiRouteService** - Configure AI routing rules
  - Model routing
  - Fallback configuration

#### Consumer Management

- **ConsumerService** - Manage API consumers
  - Consumer CRUD operations
  - Credential management
  - KeyAuth credential support (Bearer, Header, Query)
  - Allow list management

#### MCP Server Management

- **McpServerService** - Manage MCP server instances
  - OpenAPI type support
  - Database type support (MySQL, PostgreSQL, SQLite, ClickHouse)
  - Direct routing support
  - Consumer management for MCP servers

#### Kubernetes Integration

- **KubernetesClientService** - Kubernetes API client
  - Multiple authentication methods:
    - kubeconfig file path
    - kubeconfig content
    - InCluster ServiceAccount
    - Controller access token
  - CRD support:
    - WasmPlugin
    - McpBridge
    - EnvoyFilter
  - Ingress operations
  - Secret operations
  - ConfigMap operations

- **KubernetesModelConverter** - Model conversion
  - Ingress ↔ Route conversion
  - ConfigMap ↔ ServiceSource conversion
  - Secret ↔ TlsCertificate conversion
  - WasmPlugin CRD ↔ Model conversion

### Technical Details

- **Go Version**: Requires Go 1.21+
- **Dependencies**:
  - k8s.io/client-go v0.28.0
  - k8s.io/api v0.28.0
  - k8s.io/apimachinery v0.28.0
  - github.com/json-iterator/go v1.1.12
  - github.com/samber/lo v1.38.1
  - gopkg.in/yaml.v3 v3.0.1

### Documentation

- README.md - Project introduction and quick start
- README_EN.md - English documentation
- docs/usage.md - Detailed usage guide
- docs/migration.md - Java SDK migration guide

### Examples

- examples/hello-sdk - Basic SDK usage
- examples/route-management - Route CRUD operations
- examples/ai-route - LLM provider management
- examples/wasm-plugin - WASM plugin configuration

### Testing

- Unit tests for all model types
- Unit tests for configuration
- Unit tests for service layer
- Integration tests for Kubernetes operations
- Test coverage: 85%+ for model layer

---

## Version History

| Version | Date | Description |
|---------|------|-------------|
| 2.2.0 | 2026-03-25 | Initial release, ported from higress-admin-sdk Java 2.2.0 |
