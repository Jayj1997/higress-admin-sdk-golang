// Package mock provides mock implementations for testing purposes
package mock

import (
	"context"

	"github.com/Jayj1997/higress-admin-sdk-golang/v2/pkg/model"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MockKubernetesClientService is a mock implementation of KubernetesClientService
type MockKubernetesClientService struct {
	// ConfigMaps
	ConfigMaps      map[string]*corev1.ConfigMap
	ConfigMapsError error

	// Secrets
	Secrets      map[string]*corev1.Secret
	SecretsError error
}

// NewMockKubernetesClientService creates a new mock Kubernetes client
func NewMockKubernetesClientService() *MockKubernetesClientService {
	return &MockKubernetesClientService{
		ConfigMaps: make(map[string]*corev1.ConfigMap),
		Secrets:    make(map[string]*corev1.Secret),
	}
}

// ListConfigMaps lists all ConfigMaps
func (m *MockKubernetesClientService) ListConfigMaps(ctx context.Context, labels map[string]string) ([]corev1.ConfigMap, error) {
	if m.ConfigMapsError != nil {
		return nil, m.ConfigMapsError
	}

	result := make([]corev1.ConfigMap, 0, len(m.ConfigMaps))
	for _, cm := range m.ConfigMaps {
		result = append(result, *cm)
	}
	return result, nil
}

// GetConfigMap gets a ConfigMap by name
func (m *MockKubernetesClientService) GetConfigMap(ctx context.Context, name string) (*corev1.ConfigMap, error) {
	if m.ConfigMapsError != nil {
		return nil, m.ConfigMapsError
	}
	return m.ConfigMaps[name], nil
}

// CreateConfigMap creates a new ConfigMap
func (m *MockKubernetesClientService) CreateConfigMap(ctx context.Context, cm *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	if m.ConfigMapsError != nil {
		return nil, m.ConfigMapsError
	}
	m.ConfigMaps[cm.Name] = cm
	return cm, nil
}

// UpdateConfigMap updates an existing ConfigMap
func (m *MockKubernetesClientService) UpdateConfigMap(ctx context.Context, cm *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	if m.ConfigMapsError != nil {
		return nil, m.ConfigMapsError
	}
	m.ConfigMaps[cm.Name] = cm
	return cm, nil
}

// DeleteConfigMap deletes a ConfigMap
func (m *MockKubernetesClientService) DeleteConfigMap(ctx context.Context, name string) error {
	if m.ConfigMapsError != nil {
		return m.ConfigMapsError
	}
	delete(m.ConfigMaps, name)
	return nil
}

// ListSecrets lists all Secrets
func (m *MockKubernetesClientService) ListSecrets(ctx context.Context, labels map[string]string) ([]corev1.Secret, error) {
	if m.SecretsError != nil {
		return nil, m.SecretsError
	}

	result := make([]corev1.Secret, 0, len(m.Secrets))
	for _, s := range m.Secrets {
		result = append(result, *s)
	}
	return result, nil
}

// GetSecret gets a Secret by name
func (m *MockKubernetesClientService) GetSecret(ctx context.Context, name string) (*corev1.Secret, error) {
	if m.SecretsError != nil {
		return nil, m.SecretsError
	}
	return m.Secrets[name], nil
}

// CreateSecret creates a new Secret
func (m *MockKubernetesClientService) CreateSecret(ctx context.Context, s *corev1.Secret) (*corev1.Secret, error) {
	if m.SecretsError != nil {
		return nil, m.SecretsError
	}
	m.Secrets[s.Name] = s
	return s, nil
}

// UpdateSecret updates an existing Secret
func (m *MockKubernetesClientService) UpdateSecret(ctx context.Context, s *corev1.Secret) (*corev1.Secret, error) {
	if m.SecretsError != nil {
		return nil, m.SecretsError
	}
	m.Secrets[s.Name] = s
	return s, nil
}

// DeleteSecret deletes a Secret
func (m *MockKubernetesClientService) DeleteSecret(ctx context.Context, name string) error {
	if m.SecretsError != nil {
		return m.SecretsError
	}
	delete(m.Secrets, name)
	return nil
}

// MockModelConverter is a mock implementation of the model converter
type MockModelConverter struct {
	Domains      map[string]*model.Domain
	DomainsError error
}

// NewMockModelConverter creates a new mock model converter
func NewMockModelConverter() *MockModelConverter {
	return &MockModelConverter{
		Domains: make(map[string]*model.Domain),
	}
}

// ConfigMapToDomain converts a ConfigMap to a Domain
func (m *MockModelConverter) ConfigMapToDomain(cm *corev1.ConfigMap) (*model.Domain, error) {
	if m.DomainsError != nil {
		return nil, m.DomainsError
	}

	// Extract domain name from ConfigMap name
	name := cm.Name
	if len(name) > 7 && name[:7] == "domain-" {
		name = name[7:]
	}

	domain := &model.Domain{
		Name: name,
	}

	// Check if we have a stored domain
	if stored, ok := m.Domains[name]; ok {
		domain = stored
	}

	return domain, nil
}

// DomainToConfigMap converts a Domain to a ConfigMap
func (m *MockModelConverter) DomainToConfigMap(domain *model.Domain) (*corev1.ConfigMap, error) {
	if m.DomainsError != nil {
		return nil, m.DomainsError
	}

	m.Domains[domain.Name] = domain
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "domain-" + domain.Name,
			Namespace: "higress-system",
		},
		Data: map[string]string{
			"domain": domain.Name,
		},
	}, nil
}

// DomainNameToConfigMapName converts a domain name to ConfigMap name
func (m *MockModelConverter) DomainNameToConfigMapName(domainName string) string {
	return "domain-" + domainName
}

// MockRouteService is a mock implementation of RouteService
type MockRouteService struct {
	Routes      []model.Route
	RoutesError error
}

// NewMockRouteService creates a new mock route service
func NewMockRouteService() *MockRouteService {
	return &MockRouteService{
		Routes: make([]model.Route, 0),
	}
}

// List lists all routes
func (m *MockRouteService) List(ctx context.Context, query *model.RoutePageQuery) (*model.PaginatedResult[model.Route], error) {
	if m.RoutesError != nil {
		return nil, m.RoutesError
	}

	// Filter by domain if specified
	routes := m.Routes
	if query != nil && query.DomainName != "" {
		filtered := make([]model.Route, 0)
		for _, r := range routes {
			for _, d := range r.Domains {
				if d == query.DomainName {
					filtered = append(filtered, r)
					break
				}
			}
		}
		routes = filtered
	}

	return model.NewPaginatedResult(routes, len(routes), 1, 10), nil
}

// Get gets a route by name
func (m *MockRouteService) Get(ctx context.Context, name string) (*model.Route, error) {
	if m.RoutesError != nil {
		return nil, m.RoutesError
	}

	for _, r := range m.Routes {
		if r.Name == name {
			return &r, nil
		}
	}
	return nil, nil
}

// Add adds a new route
func (m *MockRouteService) Add(ctx context.Context, route *model.Route) (*model.Route, error) {
	if m.RoutesError != nil {
		return nil, m.RoutesError
	}
	m.Routes = append(m.Routes, *route)
	return route, nil
}

// Update updates a route
func (m *MockRouteService) Update(ctx context.Context, route *model.Route) (*model.Route, error) {
	if m.RoutesError != nil {
		return nil, m.RoutesError
	}

	for i, r := range m.Routes {
		if r.Name == route.Name {
			m.Routes[i] = *route
			return route, nil
		}
	}

	m.Routes = append(m.Routes, *route)
	return route, nil
}

// Delete deletes a route
func (m *MockRouteService) Delete(ctx context.Context, name string) error {
	if m.RoutesError != nil {
		return m.RoutesError
	}

	for i, r := range m.Routes {
		if r.Name == name {
			m.Routes = append(m.Routes[:i], m.Routes[i+1:]...)
			break
		}
	}
	return nil
}

// MockWasmPluginInstanceService is a mock implementation of WasmPluginInstanceService
type MockWasmPluginInstanceService struct {
	Instances      []model.WasmPluginInstance
	InstancesError error
}

// NewMockWasmPluginInstanceService creates a new mock WasmPluginInstanceService
func NewMockWasmPluginInstanceService() *MockWasmPluginInstanceService {
	return &MockWasmPluginInstanceService{
		Instances: make([]model.WasmPluginInstance, 0),
	}
}

// CreateEmptyInstance creates an empty plugin instance
func (m *MockWasmPluginInstanceService) CreateEmptyInstance(ctx context.Context, pluginName string) (*model.WasmPluginInstance, error) {
	if m.InstancesError != nil {
		return nil, m.InstancesError
	}
	return &model.WasmPluginInstance{
		PluginName: pluginName,
	}, nil
}

// ListByPlugin lists instances by plugin name
func (m *MockWasmPluginInstanceService) ListByPlugin(ctx context.Context, pluginName string, internal *bool) ([]model.WasmPluginInstance, error) {
	if m.InstancesError != nil {
		return nil, m.InstancesError
	}

	result := make([]model.WasmPluginInstance, 0)
	for _, inst := range m.Instances {
		if inst.PluginName == pluginName {
			result = append(result, inst)
		}
	}
	return result, nil
}

// ListByScope lists instances by scope
func (m *MockWasmPluginInstanceService) ListByScope(ctx context.Context, scope model.WasmPluginInstanceScope, target string) ([]model.WasmPluginInstance, error) {
	if m.InstancesError != nil {
		return nil, m.InstancesError
	}

	result := make([]model.WasmPluginInstance, 0)
	for _, inst := range m.Instances {
		if inst.Scope == scope && inst.Target == target {
			result = append(result, inst)
		}
	}
	return result, nil
}

// Query queries a specific plugin instance
func (m *MockWasmPluginInstanceService) Query(ctx context.Context, scope model.WasmPluginInstanceScope, target, pluginName string, internal *bool) (*model.WasmPluginInstance, error) {
	if m.InstancesError != nil {
		return nil, m.InstancesError
	}

	for _, inst := range m.Instances {
		if inst.Scope == scope && inst.Target == target && inst.PluginName == pluginName {
			return &inst, nil
		}
	}
	return nil, nil
}

// AddOrUpdate adds or updates an instance
func (m *MockWasmPluginInstanceService) AddOrUpdate(ctx context.Context, instance *model.WasmPluginInstance) (*model.WasmPluginInstance, error) {
	if m.InstancesError != nil {
		return nil, m.InstancesError
	}

	for i, inst := range m.Instances {
		if inst.Scope == instance.Scope && inst.Target == instance.Target && inst.PluginName == instance.PluginName {
			m.Instances[i] = *instance
			return instance, nil
		}
	}
	m.Instances = append(m.Instances, *instance)
	return instance, nil
}

// Delete deletes an instance
func (m *MockWasmPluginInstanceService) Delete(ctx context.Context, scope model.WasmPluginInstanceScope, target, pluginName string, internal *bool) error {
	if m.InstancesError != nil {
		return m.InstancesError
	}

	for i, inst := range m.Instances {
		if inst.Scope == scope && inst.Target == target && inst.PluginName == pluginName {
			m.Instances = append(m.Instances[:i], m.Instances[i+1:]...)
			break
		}
	}
	return nil
}

// DeleteAll deletes all plugin instances for a scope and target
func (m *MockWasmPluginInstanceService) DeleteAll(ctx context.Context, scope model.WasmPluginInstanceScope, target string) error {
	if m.InstancesError != nil {
		return m.InstancesError
	}

	newInstances := make([]model.WasmPluginInstance, 0)
	for _, inst := range m.Instances {
		if !(inst.Scope == scope && inst.Target == target) {
			newInstances = append(newInstances, inst)
		}
	}
	m.Instances = newInstances
	return nil
}
