// Package model provides data models for Higress Admin SDK.
package model

// WasmPlugin represents a WASM plugin definition.
type WasmPlugin struct {
	// Name is the plugin name.
	Name string `json:"name,omitempty"`

	// Version is the plugin version.
	Version string `json:"version,omitempty"`

	// PluginVersion is the plugin version (alias for Version).
	PluginVersion string `json:"pluginVersion,omitempty"`

	// Category is the plugin category.
	Category string `json:"category,omitempty"`

	// Title is the plugin display title.
	Title string `json:"title,omitempty"`

	// Description is the plugin description.
	Description string `json:"description,omitempty"`

	// ImageURL is the URL of the plugin image.
	ImageURL string `json:"imageUrl,omitempty"`

	// ImageRepository is the repository part of the image URL.
	ImageRepository string `json:"imageRepository,omitempty"`

	// ImageVersion is the version part of the image URL.
	ImageVersion string `json:"imageVersion,omitempty"`

	// Icon is the plugin icon URL or base64 encoded data.
	Icon string `json:"icon,omitempty"`

	// BuiltIn indicates whether this is a built-in plugin.
	BuiltIn *bool `json:"builtIn,omitempty"`

	// Internal indicates whether this is an internal plugin.
	Internal *bool `json:"internal,omitempty"`

	// Phase is the plugin execution phase.
	Phase string `json:"phase,omitempty"`

	// Priority is the plugin execution priority.
	Priority *int `json:"priority,omitempty"`

	// ConfigSchema is the plugin configuration schema.
	ConfigSchema map[string]interface{} `json:"configSchema,omitempty"`

	// RouteConfigSchema is the route-level configuration schema.
	RouteConfigSchema map[string]interface{} `json:"routeConfigSchema,omitempty"`

	// ImagePullPolicy is the image pull policy.
	ImagePullPolicy string `json:"imagePullPolicy,omitempty"`

	// ImagePullSecret is the image pull secret.
	ImagePullSecret string `json:"imagePullSecret,omitempty"`

	// Lang is the language for i18n.
	Lang string `json:"lang,omitempty"`
}

// WasmPluginConfig represents the configuration schema of a WASM plugin.
type WasmPluginConfig struct {
	// Schema is the OpenAPI v3 schema.
	Schema map[string]interface{} `json:"schema,omitempty"`
}

// WasmPluginInstance represents an instance of a WASM plugin.
type WasmPluginInstance struct {
	// ID is the unique identifier of the plugin instance.
	ID string `json:"id,omitempty"`

	// PluginName is the name of the plugin.
	PluginName string `json:"pluginName,omitempty"`

	// PluginVersion is the version of the plugin.
	PluginVersion string `json:"pluginVersion,omitempty"`

	// Scope is the scope of the plugin instance.
	// Valid values: "global", "domain", "route", "service"
	Scope WasmPluginInstanceScope `json:"scope,omitempty"`

	// Target is the target resource name (for non-global scopes).
	Target string `json:"target,omitempty"`

	// Targets is a map of scope to target name.
	Targets map[WasmPluginInstanceScope]string `json:"targets,omitempty"`

	// Enabled indicates whether the plugin instance is enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// Configurations are the plugin configurations.
	Configurations map[string]interface{} `json:"configurations,omitempty"`

	// Internal indicates whether this is an internal plugin instance.
	Internal *bool `json:"internal,omitempty"`
}

// WasmPluginInstanceScope represents the scope of a WASM plugin instance.
type WasmPluginInstanceScope string

// WasmPluginInstanceScope constants.
const (
	WasmPluginInstanceScopeGlobal  WasmPluginInstanceScope = "global"
	WasmPluginInstanceScopeDomain  WasmPluginInstanceScope = "domain"
	WasmPluginInstanceScopeRoute   WasmPluginInstanceScope = "route"
	WasmPluginInstanceScopeService WasmPluginInstanceScope = "service"
)

// Priority returns the priority of the scope.
func (s WasmPluginInstanceScope) Priority() int {
	switch s {
	case WasmPluginInstanceScopeGlobal:
		return 0
	case WasmPluginInstanceScopeDomain:
		return 10
	case WasmPluginInstanceScopeRoute:
		return 100
	case WasmPluginInstanceScopeService:
		return 1000
	default:
		return 0
	}
}

// Validate validates the WasmPluginInstance.
func (i *WasmPluginInstance) Validate() error {
	if i.PluginName == "" {
		return &ValidationError{Field: "pluginName", Message: "plugin name is required"}
	}
	if len(i.Targets) == 0 && i.Scope == "" {
		return &ValidationError{Field: "scope", Message: "scope or targets is required"}
	}
	return nil
}

// SyncDeprecatedFields syncs deprecated fields to the new fields.
func (i *WasmPluginInstance) SyncDeprecatedFields() {
	// Sync scope/target to targets
	if i.Scope != "" && i.Target != "" {
		if i.Targets == nil {
			i.Targets = make(map[WasmPluginInstanceScope]string)
		}
		i.Targets[i.Scope] = i.Target
	}

	// Sync targets to scope/target
	if len(i.Targets) == 1 {
		for scope, target := range i.Targets {
			i.Scope = scope
			i.Target = target
			break
		}
	}
}

// ValidationError represents a validation error.
type ValidationError struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message,omitempty"`
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	if e.Field != "" {
		return e.Field + ": " + e.Message
	}
	return e.Message
}
