// Package client provides the main client for Higress Admin SDK.
package client

import (
	"github.com/Jayj1997/higress-admin-sdk-golang/internal/kubernetes"
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/config"
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/service"
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/service/mock"
)

// HigressServiceProvider 是Higress服务提供者的主接口
// 它提供了访问所有Higress管理服务的统一入口
type HigressServiceProvider interface {
	// KubernetesClientService 返回Kubernetes客户端服务
	KubernetesClientService() *kubernetes.KubernetesClientService

	// KubernetesModelConverter 返回Kubernetes模型转换器
	KubernetesModelConverter() *kubernetes.KubernetesModelConverter

	// DomainService 返回域名管理服务
	DomainService() service.DomainService

	// RouteService 返回路由管理服务
	RouteService() service.RouteService

	// ServiceService 返回服务管理服务
	ServiceService() service.ServiceService

	// ServiceSourceService 返回服务来源管理服务
	ServiceSourceService() service.ServiceSourceService

	// ProxyServerService 返回代理服务器管理服务
	ProxyServerService() service.ProxyServerService

	// TlsCertificateService 返回TLS证书管理服务
	TlsCertificateService() service.TlsCertificateService

	// WasmPluginService 返回WASM插件管理服务
	WasmPluginService() service.WasmPluginService

	// WasmPluginInstanceService 返回WASM插件实例管理服务
	WasmPluginInstanceService() service.WasmPluginInstanceService

	// ConsumerService 返回消费者管理服务
	ConsumerService() service.ConsumerService

	// AiRouteService 返回AI路由管理服务
	AiRouteService() service.AiRouteService

	// LlmProviderService 返回LLM提供商管理服务
	LlmProviderService() service.LlmProviderService

	// McpServerService 返回MCP服务器管理服务
	McpServerService() service.McpServerService
}

// HigressServiceProviderImpl HigressServiceProvider的实现
type HigressServiceProviderImpl struct {
	kubernetesClientService   *kubernetes.KubernetesClientService
	kubernetesModelConverter  *kubernetes.KubernetesModelConverter
	domainService             service.DomainService
	routeService              service.RouteService
	serviceService            service.ServiceService
	serviceSourceService      service.ServiceSourceService
	proxyServerService        service.ProxyServerService
	tlsCertificateService     service.TlsCertificateService
	wasmPluginService         service.WasmPluginService
	wasmPluginInstanceService service.WasmPluginInstanceService
	consumerService           service.ConsumerService
	aiRouteService            service.AiRouteService
	llmProviderService        service.LlmProviderService
	mcpServerService          service.McpServerService
}

// NewHigressServiceProvider 创建新的HigressServiceProvider实例
func NewHigressServiceProvider(cfg *config.HigressServiceConfig) (HigressServiceProvider, error) {
	// 1. 创建Kubernetes客户端服务
	kubernetesClientService, err := kubernetes.NewKubernetesClientService(cfg)
	if err != nil {
		return nil, err
	}

	// 2. 创建Kubernetes模型转换器
	kubernetesModelConverter := kubernetes.NewKubernetesModelConverter(kubernetesClientService)

	// 3. 创建ServiceService
	serviceService := service.NewServiceService(kubernetesClientService)

	// 4. 创建ServiceSourceService
	serviceSourceService := service.NewServiceSourceService(kubernetesClientService, kubernetesModelConverter)

	// 5. 创建ProxyServerService
	proxyServerService := service.NewProxyServerService(kubernetesClientService, kubernetesModelConverter)

	// 6. 创建TlsCertificateService
	tlsCertificateService := service.NewTlsCertificateService(kubernetesClientService, kubernetesModelConverter)

	// 7. 创建WasmPluginService (Mock实现，将在里程碑7中完整实现)
	wasmPluginService := mock.NewMockWasmPluginService()

	// 8. 创建WasmPluginInstanceService (Mock实现，将在里程碑7中完整实现)
	wasmPluginInstanceService := service.NewMockWasmPluginInstanceService()

	// 9. 创建ConsumerService (Mock实现，将在里程碑9中完整实现)
	consumerService := mock.NewMockConsumerService()

	// 10. 创建RouteService
	routeService := service.NewRouteService(kubernetesClientService, kubernetesModelConverter, wasmPluginInstanceService)

	// 11. 创建DomainService
	domainService := service.NewDomainService(kubernetesClientService, kubernetesModelConverter, routeService, wasmPluginInstanceService)

	// 12. 创建LlmProviderService (Mock实现，将在里程碑8中完整实现)
	llmProviderService := mock.NewMockLlmProviderService()

	// 13. 创建AiRouteService (Mock实现，将在里程碑8中完整实现)
	aiRouteService := mock.NewMockAiRouteService()

	// 14. 创建McpServerService (Mock实现，将在里程碑10中完整实现)
	mcpServerService := mock.NewMockMcpServerService()

	return &HigressServiceProviderImpl{
		kubernetesClientService:   kubernetesClientService,
		kubernetesModelConverter:  kubernetesModelConverter,
		domainService:             domainService,
		routeService:              routeService,
		serviceService:            serviceService,
		serviceSourceService:      serviceSourceService,
		proxyServerService:        proxyServerService,
		tlsCertificateService:     tlsCertificateService,
		wasmPluginService:         wasmPluginService,
		wasmPluginInstanceService: wasmPluginInstanceService,
		consumerService:           consumerService,
		aiRouteService:            aiRouteService,
		llmProviderService:        llmProviderService,
		mcpServerService:          mcpServerService,
	}, nil
}

// KubernetesClientService 返回Kubernetes客户端服务
func (p *HigressServiceProviderImpl) KubernetesClientService() *kubernetes.KubernetesClientService {
	return p.kubernetesClientService
}

// KubernetesModelConverter 返回Kubernetes模型转换器
func (p *HigressServiceProviderImpl) KubernetesModelConverter() *kubernetes.KubernetesModelConverter {
	return p.kubernetesModelConverter
}

// DomainService 返回域名管理服务
func (p *HigressServiceProviderImpl) DomainService() service.DomainService {
	return p.domainService
}

// RouteService 返回路由管理服务
func (p *HigressServiceProviderImpl) RouteService() service.RouteService {
	return p.routeService
}

// ServiceService 返回服务管理服务
func (p *HigressServiceProviderImpl) ServiceService() service.ServiceService {
	return p.serviceService
}

// ServiceSourceService 返回服务来源管理服务
func (p *HigressServiceProviderImpl) ServiceSourceService() service.ServiceSourceService {
	return p.serviceSourceService
}

// ProxyServerService 返回代理服务器管理服务
func (p *HigressServiceProviderImpl) ProxyServerService() service.ProxyServerService {
	return p.proxyServerService
}

// TlsCertificateService 返回TLS证书管理服务
func (p *HigressServiceProviderImpl) TlsCertificateService() service.TlsCertificateService {
	return p.tlsCertificateService
}

// WasmPluginService 返回WASM插件管理服务
func (p *HigressServiceProviderImpl) WasmPluginService() service.WasmPluginService {
	return p.wasmPluginService
}

// WasmPluginInstanceService 返回WASM插件实例管理服务
func (p *HigressServiceProviderImpl) WasmPluginInstanceService() service.WasmPluginInstanceService {
	return p.wasmPluginInstanceService
}

// ConsumerService 返回消费者管理服务
func (p *HigressServiceProviderImpl) ConsumerService() service.ConsumerService {
	return p.consumerService
}

// AiRouteService 返回AI路由管理服务
func (p *HigressServiceProviderImpl) AiRouteService() service.AiRouteService {
	return p.aiRouteService
}

// LlmProviderService 返回LLM提供商管理服务
func (p *HigressServiceProviderImpl) LlmProviderService() service.LlmProviderService {
	return p.llmProviderService
}

// McpServerService 返回MCP服务器管理服务
func (p *HigressServiceProviderImpl) McpServerService() service.McpServerService {
	return p.mcpServerService
}
