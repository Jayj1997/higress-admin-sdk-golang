// Package model provides data models for Higress Admin SDK.
package model

// Plugin represents the plugin spec.yaml structure.
type Plugin struct {
	// APIVersion is the API version.
	APIVersion string `yaml:"apiVersion"`

	// Info contains the plugin information.
	Info PluginInfo `yaml:"info"`

	// Spec contains the plugin specification.
	Spec PluginSpec `yaml:"spec"`
}

// PluginInfo contains basic plugin information.
type PluginInfo struct {
	// GatewayMinVersion is the minimum gateway version required.
	GatewayMinVersion string `yaml:"gatewayMinVersion"`

	// Type is the plugin type (e.g., enterprise).
	Type string `yaml:"type"`

	// Category is the plugin category (e.g., auth, ai, traffic).
	Category string `yaml:"category"`

	// Name is the plugin name.
	Name string `yaml:"name"`

	// Image is the plugin image path.
	Image string `yaml:"image"`

	// Title is the plugin display title.
	Title string `yaml:"title"`

	// TitleI18n contains internationalized titles.
	TitleI18n map[string]string `yaml:"x-title-i18n"`

	// Description is the plugin description.
	Description string `yaml:"description"`

	// DescriptionI18n contains internationalized descriptions.
	DescriptionI18n map[string]string `yaml:"x-description-i18n"`

	// IconURL is the URL to the plugin icon.
	IconURL string `yaml:"iconUrl"`

	// Version is the plugin version.
	Version string `yaml:"version"`

	// Contact contains contact information.
	Contact *PluginContact `yaml:"contact"`
}

// PluginContact contains contact information.
type PluginContact struct {
	Name  string `yaml:"name"`
	URL   string `yaml:"url"`
	Email string `yaml:"email"`
}

// PluginSpec contains the plugin specification.
type PluginSpec struct {
	// Phase is the execution phase (e.g., AUTHN, AUTHZ, STATS).
	Phase string `yaml:"phase"`

	// Priority is the execution priority.
	Priority int `yaml:"priority"`

	// ConfigSchema is the global configuration schema.
	ConfigSchema *ConfigSchema `yaml:"configSchema"`

	// RouteConfigSchema is the route-level configuration schema.
	RouteConfigSchema *ConfigSchema `yaml:"routeConfigSchema"`
}

// ConfigSchema represents the OpenAPI v3 schema for plugin configuration.
type ConfigSchema struct {
	// OpenAPIV3Schema is the OpenAPI v3 schema definition.
	OpenAPIV3Schema map[string]interface{} `yaml:"openAPIV3Schema"`
}

// GetTitle returns the title for the specified language.
func (i *PluginInfo) GetTitle(lang string) string {
	if lang != "" && i.TitleI18n != nil {
		if title, ok := i.TitleI18n[lang]; ok {
			return title
		}
	}
	return i.Title
}

// GetDescription returns the description for the specified language.
func (i *PluginInfo) GetDescription(lang string) string {
	if lang != "" && i.DescriptionI18n != nil {
		if desc, ok := i.DescriptionI18n[lang]; ok {
			return desc
		}
	}
	return i.Description
}

// ToWasmPlugin converts Plugin to WasmPlugin model.
func (p *Plugin) ToWasmPlugin(lang string) *WasmPlugin {
	builtIn := true
	return &WasmPlugin{
		Name:        p.Info.Name,
		Version:     p.Info.Version,
		Category:    p.Info.Category,
		Title:       p.Info.GetTitle(lang),
		Description: p.Info.GetDescription(lang),
		Icon:        p.Info.IconURL,
		BuiltIn:     &builtIn,
		Phase:       p.Spec.Phase,
		Priority:    &p.Spec.Priority,
	}
}
