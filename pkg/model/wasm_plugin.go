// Package model provides data models for Higress Admin SDK.
package model

// WasmPlugin represents a WASM plugin definition.
type WasmPlugin struct {
	// Name is the plugin name.
	Name string `json:"name,omitempty"`

	// Version is the plugin version.
	Version string `json:"version,omitempty"`

	// Category is the plugin category.
	Category string `json:"category,omitempty"`

	// Title is the plugin display title.
	Title string `json:"title,omitempty"`

	// Description is the plugin description.
	Description string `json:"description,omitempty"`

	// ImageURL is the URL of the plugin image.
	ImageURL string `json:"imageUrl,omitempty"`

	// Icon is the plugin icon URL.
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
