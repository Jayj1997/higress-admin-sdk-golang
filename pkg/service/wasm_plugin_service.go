// Package service provides business services for the SDK
package service

import (
	"context"
	"encoding/base64"
	"sort"
	"strings"
	"sync"

	"github.com/Jayj1997/higress-admin-sdk-golang/internal/kubernetes"
	"github.com/Jayj1997/higress-admin-sdk-golang/internal/resources/plugins"
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/config"
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/errors"
	"github.com/Jayj1997/higress-admin-sdk-golang/pkg/model"
	"gopkg.in/yaml.v3"
)

// WasmPluginService WASM插件管理服务接口
type WasmPluginService interface {
	// List 列出WASM插件（内置+自定义）
	List(ctx context.Context, query *model.WasmPluginPageQuery) (*model.PaginatedResult[model.WasmPlugin], error)

	// Get 获取WASM插件详情
	Get(ctx context.Context, name, language string) (*model.WasmPlugin, error)

	// GetConfig 获取插件配置Schema
	GetConfig(ctx context.Context, name, language string) (*model.WasmPluginConfig, error)

	// GetReadme 获取插件README文档
	GetReadme(ctx context.Context, name, language string) (string, error)

	// UpdateBuiltIn 更新内置插件配置
	UpdateBuiltIn(ctx context.Context, plugin *model.WasmPlugin) (*model.WasmPlugin, error)

	// AddCustom 添加自定义插件
	AddCustom(ctx context.Context, plugin *model.WasmPlugin) (*model.WasmPlugin, error)

	// UpdateCustom 更新自定义插件
	UpdateCustom(ctx context.Context, plugin *model.WasmPlugin) (*model.WasmPlugin, error)

	// DeleteCustom 删除自定义插件
	DeleteCustom(ctx context.Context, name string) error
}

// WasmPluginServiceImpl implements WasmPluginService
type WasmPluginServiceImpl struct {
	kubernetesClient *kubernetes.KubernetesClientService
	modelConverter   *kubernetes.KubernetesModelConverter
	config           *config.WasmPluginServiceConfig
	builtInPlugins   []*pluginCacheItem
	builtInPluginsMu sync.RWMutex
}

// pluginCacheItem represents a cached built-in plugin
type pluginCacheItem struct {
	name            string
	plugin          *model.Plugin
	imageURL        string
	imagePullSecret string
	imagePullPolicy string
	readme          map[string]string // lang -> content
	readmeDefault   string
	iconData        string
}

// NewWasmPluginService creates a new WasmPluginService
func NewWasmPluginService(
	kubernetesClient *kubernetes.KubernetesClientService,
	modelConverter *kubernetes.KubernetesModelConverter,
	cfg *config.WasmPluginServiceConfig,
) (WasmPluginService, error) {
	svc := &WasmPluginServiceImpl{
		kubernetesClient: kubernetesClient,
		modelConverter:   modelConverter,
		config:           cfg,
	}
	if err := svc.initialize(); err != nil {
		return nil, err
	}
	return svc, nil
}

// initialize loads all built-in plugins
func (s *WasmPluginServiceImpl) initialize() error {
	pluginNames := plugins.ListPlugins()
	items := make([]*pluginCacheItem, 0, len(pluginNames))

	for _, name := range pluginNames {
		item, err := s.loadPlugin(name)
		if err != nil {
			// Log error but continue loading other plugins
			continue
		}
		if item != nil {
			items = append(items, item)
		}
	}

	// Sort by name
	sort.Slice(items, func(i, j int) bool {
		return items[i].name < items[j].name
	})

	s.builtInPluginsMu.Lock()
	s.builtInPlugins = items
	s.builtInPluginsMu.Unlock()

	return nil
}

// loadPlugin loads a single plugin from embedded resources
func (s *WasmPluginServiceImpl) loadPlugin(name string) (*pluginCacheItem, error) {
	// Load spec.yaml
	specData, err := plugins.GetPluginSpec(name)
	if err != nil {
		return nil, err
	}

	var plugin model.Plugin
	if err := yaml.Unmarshal(specData, &plugin); err != nil {
		return nil, err
	}

	// Extract config example from spec
	s.fillPluginConfigExample(&plugin, string(specData))

	item := &pluginCacheItem{
		name:   name,
		plugin: &plugin,
		readme: make(map[string]string),
	}

	// Get image URL
	if s.config != nil && s.config.CustomImageUrlPattern != "" {
		item.imageURL = s.formatImageUrl(s.config.CustomImageUrlPattern, &plugin.Info)
	} else {
		item.imageURL = plugins.GetPluginImageURL(name)
	}

	if s.config != nil {
		item.imagePullSecret = s.config.ImagePullSecret
		item.imagePullPolicy = s.config.ImagePullPolicy
	}

	// Load README files
	if content, err := plugins.GetPluginReadme(name, ""); err == nil {
		item.readmeDefault = string(content)
	}
	if content, err := plugins.GetPluginReadme(name, "zh-CN"); err == nil {
		item.readme["zh-CN"] = string(content)
	}
	if content, err := plugins.GetPluginReadme(name, "en-US"); err == nil {
		item.readme["en-US"] = string(content)
	}

	// Load icon
	if iconData, err := plugins.GetPluginIcon(name); err == nil {
		item.iconData = "data:image/png;base64," + base64.StdEncoding.EncodeToString(iconData)
	}

	return item, nil
}

// formatImageUrl formats the image URL using the pattern
func (s *WasmPluginServiceImpl) formatImageUrl(pattern string, info *model.PluginInfo) string {
	result := pattern
	result = strings.ReplaceAll(result, "${name}", info.Name)
	result = strings.ReplaceAll(result, "${version}", info.Version)
	return result
}

// fillPluginConfigExample extracts and fills the config example from spec
func (s *WasmPluginServiceImpl) fillPluginConfigExample(plugin *model.Plugin, content string) {
	example := s.extractConfigExample(content)
	if example == "" {
		return
	}
	if plugin.Spec.ConfigSchema != nil && plugin.Spec.ConfigSchema.OpenAPIV3Schema != nil {
		if plugin.Spec.ConfigSchema.OpenAPIV3Schema == nil {
			plugin.Spec.ConfigSchema.OpenAPIV3Schema = make(map[string]interface{})
		}
		plugin.Spec.ConfigSchema.OpenAPIV3Schema["x-example-raw"] = example
	}
}

// extractConfigExample extracts the example section from spec.yaml
func (s *WasmPluginServiceImpl) extractConfigExample(content string) string {
	lines := strings.Split(content, "\n")
	var result strings.Builder
	var inConfigSchema, inOpenApiV3Schema, inExample bool
	var schemaIndent, exampleIndent string

	for _, line := range lines {
		trimmed := strings.TrimLeft(line, " \t")
		indent := line[:len(line)-len(trimmed)]

		if !inConfigSchema {
			if strings.HasPrefix(trimmed, "configSchema:") {
				inConfigSchema = true
				schemaIndent = indent
			}
			continue
		}

		if !inOpenApiV3Schema {
			if strings.HasPrefix(trimmed, "openAPIV3Schema:") {
				inOpenApiV3Schema = true
				schemaIndent = indent
			}
			continue
		}

		// Check if we've exited the schema block
		if len(indent) <= len(schemaIndent) && !strings.HasPrefix(trimmed, "openAPIV3Schema:") {
			break
		}

		if !inExample {
			if strings.HasPrefix(trimmed, "example:") {
				inExample = true
				exampleIndent = indent
			}
			continue
		}

		// Check if we've exited the example block
		if len(indent) <= len(exampleIndent) {
			break
		}

		// Remove one level of indentation
		if strings.HasPrefix(line, exampleIndent+"  ") {
			result.WriteString(line[len(exampleIndent)+2:])
		} else if strings.HasPrefix(line, exampleIndent+"\t") {
			result.WriteString(line[len(exampleIndent)+1:])
		} else {
			result.WriteString(trimmed)
		}
		result.WriteString("\n")
	}

	return strings.TrimSpace(result.String())
}

// List lists all plugins (built-in + custom)
func (s *WasmPluginServiceImpl) List(ctx context.Context, query *model.WasmPluginPageQuery) (*model.PaginatedResult[model.WasmPlugin], error) {
	lang := ""
	if query != nil {
		lang = query.Lang
	}

	var result []model.WasmPlugin

	// Add built-in plugins
	s.builtInPluginsMu.RLock()
	for _, item := range s.builtInPlugins {
		plugin := item.buildWasmPlugin(lang)
		result = append(result, *plugin)
	}
	s.builtInPluginsMu.RUnlock()

	// Add custom plugins from Kubernetes
	if s.kubernetesClient != nil {
		crs, err := s.kubernetesClient.ListWasmPlugins(ctx, "", "", nil)
		if err != nil {
			return nil, err
		}
		for _, cr := range crs {
			plugin, err := s.modelConverter.WasmPluginCRDToModel(cr)
			if err != nil {
				continue
			}
			if plugin.BuiltIn != nil && *plugin.BuiltIn {
				// Update existing built-in plugin with K8s data
				for i, p := range result {
					if p.Name == plugin.Name {
						result[i].ImageRepository = plugin.ImageRepository
						result[i].ImageVersion = plugin.ImageVersion
						result[i].Phase = plugin.Phase
						result[i].Priority = plugin.Priority
						result[i].ImagePullPolicy = plugin.ImagePullPolicy
						result[i].ImagePullSecret = plugin.ImagePullSecret
						break
					}
				}
			} else {
				result = append(result, *plugin)
			}
		}
	}

	// Apply filters
	if query != nil {
		if query.Name != "" {
			var filtered []model.WasmPlugin
			for _, p := range result {
				if p.Name == query.Name {
					filtered = append(filtered, p)
				}
			}
			result = filtered
		}
		if query.BuiltIn != nil {
			var filtered []model.WasmPlugin
			for _, p := range result {
				if p.BuiltIn != nil && *p.BuiltIn == *query.BuiltIn {
					filtered = append(filtered, p)
				}
			}
			result = filtered
		}
		if query.Category != "" {
			var filtered []model.WasmPlugin
			for _, p := range result {
				if p.Category == query.Category {
					filtered = append(filtered, p)
				}
			}
			result = filtered
		}
	}

	// Apply pagination
	total := len(result)
	pageNum := 1
	pageSize := total
	if query != nil {
		if query.PageNum > 0 {
			pageNum = query.PageNum
		}
		if query.PageSize > 0 {
			pageSize = query.PageSize
		}
	}

	// Apply pagination
	start := (pageNum - 1) * pageSize
	end := start + pageSize
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}
	pagedResult := result[start:end]

	return model.NewPaginatedResult(pagedResult, total, pageNum, pageSize), nil
}

// Get gets a plugin by name
func (s *WasmPluginServiceImpl) Get(ctx context.Context, name, language string) (*model.WasmPlugin, error) {
	if name == "" {
		return nil, nil
	}

	// Check built-in plugins
	s.builtInPluginsMu.RLock()
	for _, item := range s.builtInPlugins {
		if item.name == name {
			result := item.buildWasmPlugin(language)
			s.builtInPluginsMu.RUnlock()
			return result, nil
		}
	}
	s.builtInPluginsMu.RUnlock()

	// Check custom plugins
	if s.kubernetesClient != nil {
		crs, err := s.kubernetesClient.ListWasmPlugins(ctx, name, "", nil)
		if err != nil {
			return nil, err
		}
		if len(crs) > 0 {
			// Return the one with highest priority
			maxPriority := 0
			var result *model.WasmPlugin
			for _, cr := range crs {
				plugin, err := s.modelConverter.WasmPluginCRDToModel(cr)
				if err != nil {
					continue
				}
				priority := 0
				if plugin.Priority != nil {
					priority = *plugin.Priority
				}
				if priority >= maxPriority {
					maxPriority = priority
					result = plugin
				}
			}
			return result, nil
		}
	}

	return nil, nil
}

// GetConfig gets the plugin configuration schema
func (s *WasmPluginServiceImpl) GetConfig(ctx context.Context, name, language string) (*model.WasmPluginConfig, error) {
	if name == "" {
		return nil, nil
	}

	// Check built-in plugins
	s.builtInPluginsMu.RLock()
	for _, item := range s.builtInPlugins {
		if item.name == name {
			config := item.buildWasmPluginConfig(language)
			s.builtInPluginsMu.RUnlock()
			return config, nil
		}
	}
	s.builtInPluginsMu.RUnlock()

	// For custom plugins, return an empty schema
	if s.kubernetesClient != nil {
		crs, err := s.kubernetesClient.ListWasmPlugins(ctx, name, "", nil)
		if err != nil {
			return nil, err
		}
		if len(crs) > 0 {
			return &model.WasmPluginConfig{
				Schema: map[string]interface{}{
					"type": "object",
				},
			}, nil
		}
	}

	return nil, nil
}

// GetReadme gets the plugin README
func (s *WasmPluginServiceImpl) GetReadme(ctx context.Context, name, language string) (string, error) {
	if name == "" {
		return "", nil
	}

	// Check built-in plugins
	s.builtInPluginsMu.RLock()
	for _, item := range s.builtInPlugins {
		if item.name == name {
			var content string
			if language != "" {
				content = item.readme[language]
			}
			if content == "" {
				content = item.readmeDefault
			}
			s.builtInPluginsMu.RUnlock()
			return content, nil
		}
	}
	s.builtInPluginsMu.RUnlock()

	return "", nil
}

// UpdateBuiltIn updates a built-in plugin configuration
func (s *WasmPluginServiceImpl) UpdateBuiltIn(ctx context.Context, plugin *model.WasmPlugin) (*model.WasmPlugin, error) {
	// Built-in plugins are updated by creating/updating the WasmPlugin CR
	return s.AddCustom(ctx, plugin)
}

// AddCustom adds a custom plugin
func (s *WasmPluginServiceImpl) AddCustom(ctx context.Context, plugin *model.WasmPlugin) (*model.WasmPlugin, error) {
	if s.kubernetesClient == nil {
		return nil, errors.NewBusinessError("kubernetes client is not available")
	}

	cr, err := s.modelConverter.ModelToWasmPluginCRD(plugin)
	if err != nil {
		return nil, err
	}
	created, err := s.kubernetesClient.CreateWasmPlugin(ctx, cr)
	if err != nil {
		return nil, err
	}

	return s.modelConverter.WasmPluginCRDToModel(created)
}

// UpdateCustom updates a custom plugin
func (s *WasmPluginServiceImpl) UpdateCustom(ctx context.Context, plugin *model.WasmPlugin) (*model.WasmPlugin, error) {
	if s.kubernetesClient == nil {
		return nil, errors.NewBusinessError("kubernetes client is not available")
	}

	cr, err := s.modelConverter.ModelToWasmPluginCRD(plugin)
	if err != nil {
		return nil, err
	}
	updated, err := s.kubernetesClient.UpdateWasmPlugin(ctx, cr)
	if err != nil {
		return nil, err
	}

	return s.modelConverter.WasmPluginCRDToModel(updated)
}

// DeleteCustom deletes a custom plugin
func (s *WasmPluginServiceImpl) DeleteCustom(ctx context.Context, name string) error {
	if s.kubernetesClient == nil {
		return errors.NewBusinessError("kubernetes client is not available")
	}

	return s.kubernetesClient.DeleteWasmPlugin(ctx, name)
}

// buildWasmPlugin builds a WasmPlugin from the cache item
func (item *pluginCacheItem) buildWasmPlugin(lang string) *model.WasmPlugin {
	builtIn := true
	plugin := &model.WasmPlugin{
		Name:        item.plugin.Info.Name,
		Version:     item.plugin.Info.Version,
		Category:    item.plugin.Info.Category,
		Title:       item.plugin.Info.GetTitle(lang),
		Description: item.plugin.Info.GetDescription(lang),
		ImageURL:    item.imageURL,
		Icon:        item.iconData,
		BuiltIn:     &builtIn,
		Phase:       item.plugin.Spec.Phase,
		Priority:    &item.plugin.Spec.Priority,
	}

	if item.plugin.Spec.ConfigSchema != nil && item.plugin.Spec.ConfigSchema.OpenAPIV3Schema != nil {
		plugin.ConfigSchema = item.plugin.Spec.ConfigSchema.OpenAPIV3Schema
	}
	if item.plugin.Spec.RouteConfigSchema != nil && item.plugin.Spec.RouteConfigSchema.OpenAPIV3Schema != nil {
		plugin.RouteConfigSchema = item.plugin.Spec.RouteConfigSchema.OpenAPIV3Schema
	}

	return plugin
}

// buildWasmPluginConfig builds a WasmPluginConfig from the cache item
func (item *pluginCacheItem) buildWasmPluginConfig(lang string) *model.WasmPluginConfig {
	if item.plugin.Spec.ConfigSchema == nil || item.plugin.Spec.ConfigSchema.OpenAPIV3Schema == nil {
		return &model.WasmPluginConfig{Schema: map[string]interface{}{"type": "object"}}
	}

	// Apply i18n to schema
	schema := applyI18nToSchema(item.plugin.Spec.ConfigSchema.OpenAPIV3Schema, lang)
	return &model.WasmPluginConfig{Schema: schema}
}

// applyI18nToSchema applies i18n to the schema
func applyI18nToSchema(schema map[string]interface{}, lang string) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range schema {
		if strings.HasPrefix(k, "x-") && strings.HasSuffix(k, "-i18n") {
			// Extract i18n value
			if m, ok := v.(map[string]interface{}); ok && lang != "" {
				if val, exists := m[lang]; exists {
					// Get the base key
					baseKey := k[2 : len(k)-5]
					if _, hasBase := schema[baseKey]; hasBase {
						result[baseKey] = val
						continue
					}
				}
			}
			continue
		}
		switch val := v.(type) {
		case map[string]interface{}:
			result[k] = applyI18nToSchema(val, lang)
		case []interface{}:
			result[k] = val
		default:
			result[k] = v
		}
	}
	return result
}
