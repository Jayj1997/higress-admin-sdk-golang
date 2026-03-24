// Package model provides data models for Higress Admin SDK.
package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPluginInfo_GetTitle(t *testing.T) {
	tests := []struct {
		name     string
		info     *PluginInfo
		lang     string
		expected string
	}{
		{
			name: "returns default title when lang is empty",
			info: &PluginInfo{
				Title:       "Default Title",
				TitleI18n:   map[string]string{"zh": "中文标题"},
				Description: "Description",
			},
			lang:     "",
			expected: "Default Title",
		},
		{
			name: "returns default title when TitleI18n is nil",
			info: &PluginInfo{
				Title:       "Default Title",
				TitleI18n:   nil,
				Description: "Description",
			},
			lang:     "zh",
			expected: "Default Title",
		},
		{
			name: "returns i18n title when lang matches",
			info: &PluginInfo{
				Title:       "Default Title",
				TitleI18n:   map[string]string{"zh": "中文标题", "en": "English Title"},
				Description: "Description",
			},
			lang:     "zh",
			expected: "中文标题",
		},
		{
			name: "returns default title when lang not found in i18n",
			info: &PluginInfo{
				Title:       "Default Title",
				TitleI18n:   map[string]string{"zh": "中文标题"},
				Description: "Description",
			},
			lang:     "fr",
			expected: "Default Title",
		},
		{
			name: "returns empty title when both are empty",
			info: &PluginInfo{
				Title:       "",
				TitleI18n:   map[string]string{},
				Description: "Description",
			},
			lang:     "en",
			expected: "",
		},
		{
			name: "returns i18n title for english",
			info: &PluginInfo{
				Title:       "Default Title",
				TitleI18n:   map[string]string{"en": "English Title"},
				Description: "Description",
			},
			lang:     "en",
			expected: "English Title",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.info.GetTitle(tt.lang)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPluginInfo_GetDescription(t *testing.T) {
	tests := []struct {
		name     string
		info     *PluginInfo
		lang     string
		expected string
	}{
		{
			name: "returns default description when lang is empty",
			info: &PluginInfo{
				Title:           "Title",
				Description:     "Default Description",
				DescriptionI18n: map[string]string{"zh": "中文描述"},
			},
			lang:     "",
			expected: "Default Description",
		},
		{
			name: "returns default description when DescriptionI18n is nil",
			info: &PluginInfo{
				Title:           "Title",
				Description:     "Default Description",
				DescriptionI18n: nil,
			},
			lang:     "zh",
			expected: "Default Description",
		},
		{
			name: "returns i18n description when lang matches",
			info: &PluginInfo{
				Title:           "Title",
				Description:     "Default Description",
				DescriptionI18n: map[string]string{"zh": "中文描述", "en": "English Description"},
			},
			lang:     "zh",
			expected: "中文描述",
		},
		{
			name: "returns default description when lang not found in i18n",
			info: &PluginInfo{
				Title:           "Title",
				Description:     "Default Description",
				DescriptionI18n: map[string]string{"zh": "中文描述"},
			},
			lang:     "fr",
			expected: "Default Description",
		},
		{
			name: "returns empty description when both are empty",
			info: &PluginInfo{
				Title:           "Title",
				Description:     "",
				DescriptionI18n: map[string]string{},
			},
			lang:     "en",
			expected: "",
		},
		{
			name: "returns i18n description for english",
			info: &PluginInfo{
				Title:           "Title",
				Description:     "Default Description",
				DescriptionI18n: map[string]string{"en": "English Description"},
			},
			lang:     "en",
			expected: "English Description",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.info.GetDescription(tt.lang)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPlugin_ToWasmPlugin(t *testing.T) {
	tests := []struct {
		name     string
		plugin   *Plugin
		lang     string
		expected *WasmPlugin
	}{
		{
			name: "converts plugin with basic info",
			plugin: &Plugin{
				APIVersion: "v1",
				Info: PluginInfo{
					Name:        "basic-auth",
					Version:     "1.0.0",
					Category:    "auth",
					Title:       "Basic Auth",
					Description: "Basic authentication plugin",
					IconURL:     "https://example.com/icon.png",
				},
				Spec: PluginSpec{
					Phase:    "AUTHN",
					Priority: 100,
				},
			},
			lang: "en",
			expected: &WasmPlugin{
				Name:        "basic-auth",
				Version:     "1.0.0",
				Category:    "auth",
				Title:       "Basic Auth",
				Description: "Basic authentication plugin",
				Icon:        "https://example.com/icon.png",
				BuiltIn:     boolPtr(true),
				Phase:       "AUTHN",
				Priority:    intPtr(100),
			},
		},
		{
			name: "converts plugin with i18n title and description",
			plugin: &Plugin{
				APIVersion: "v1",
				Info: PluginInfo{
					Name:            "ai-proxy",
					Version:         "1.0.0",
					Category:        "ai",
					Title:           "AI Proxy",
					TitleI18n:       map[string]string{"zh": "AI代理"},
					Description:     "AI proxy plugin",
					DescriptionI18n: map[string]string{"zh": "AI代理插件"},
					IconURL:         "https://example.com/ai-icon.png",
				},
				Spec: PluginSpec{
					Phase:    "AUTHN",
					Priority: 200,
				},
			},
			lang: "zh",
			expected: &WasmPlugin{
				Name:        "ai-proxy",
				Version:     "1.0.0",
				Category:    "ai",
				Title:       "AI代理",
				Description: "AI代理插件",
				Icon:        "https://example.com/ai-icon.png",
				BuiltIn:     boolPtr(true),
				Phase:       "AUTHN",
				Priority:    intPtr(200),
			},
		},
		{
			name: "converts plugin with contact info",
			plugin: &Plugin{
				APIVersion: "v1",
				Info: PluginInfo{
					Name:        "custom-plugin",
					Version:     "2.0.0",
					Category:    "traffic",
					Title:       "Custom Plugin",
					Description: "A custom plugin",
					IconURL:     "https://example.com/custom.png",
					Contact: &PluginContact{
						Name:  "Developer",
						URL:   "https://example.com",
						Email: "dev@example.com",
					},
				},
				Spec: PluginSpec{
					Phase:    "STATS",
					Priority: 50,
				},
			},
			lang: "en",
			expected: &WasmPlugin{
				Name:        "custom-plugin",
				Version:     "2.0.0",
				Category:    "traffic",
				Title:       "Custom Plugin",
				Description: "A custom plugin",
				Icon:        "https://example.com/custom.png",
				BuiltIn:     boolPtr(true),
				Phase:       "STATS",
				Priority:    intPtr(50),
			},
		},
		{
			name: "converts plugin with empty lang returns default values",
			plugin: &Plugin{
				APIVersion: "v1",
				Info: PluginInfo{
					Name:            "test-plugin",
					Version:         "1.0.0",
					Category:        "test",
					Title:           "Test Plugin",
					TitleI18n:       map[string]string{"zh": "测试插件"},
					Description:     "Test description",
					DescriptionI18n: map[string]string{"zh": "测试描述"},
					IconURL:         "https://example.com/test.png",
				},
				Spec: PluginSpec{
					Phase:    "AUTHZ",
					Priority: 75,
				},
			},
			lang: "",
			expected: &WasmPlugin{
				Name:        "test-plugin",
				Version:     "1.0.0",
				Category:    "test",
				Title:       "Test Plugin",
				Description: "Test description",
				Icon:        "https://example.com/test.png",
				BuiltIn:     boolPtr(true),
				Phase:       "AUTHZ",
				Priority:    intPtr(75),
			},
		},
		{
			name: "converts plugin with zero priority",
			plugin: &Plugin{
				APIVersion: "v1",
				Info: PluginInfo{
					Name:        "zero-priority",
					Version:     "1.0.0",
					Category:    "test",
					Title:       "Zero Priority",
					Description: "Zero priority plugin",
				},
				Spec: PluginSpec{
					Phase:    "AUTHN",
					Priority: 0,
				},
			},
			lang: "en",
			expected: &WasmPlugin{
				Name:        "zero-priority",
				Version:     "1.0.0",
				Category:    "test",
				Title:       "Zero Priority",
				Description: "Zero priority plugin",
				Icon:        "",
				BuiltIn:     boolPtr(true),
				Phase:       "AUTHN",
				Priority:    intPtr(0),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.plugin.ToWasmPlugin(tt.lang)
			assert.Equal(t, tt.expected.Name, result.Name)
			assert.Equal(t, tt.expected.Version, result.Version)
			assert.Equal(t, tt.expected.Category, result.Category)
			assert.Equal(t, tt.expected.Title, result.Title)
			assert.Equal(t, tt.expected.Description, result.Description)
			assert.Equal(t, tt.expected.Icon, result.Icon)
			assert.Equal(t, tt.expected.BuiltIn, result.BuiltIn)
			assert.Equal(t, tt.expected.Phase, result.Phase)
			assert.Equal(t, tt.expected.Priority, result.Priority)
		})
	}
}

func TestPluginStructFields(t *testing.T) {
	t.Run("Plugin struct with all fields", func(t *testing.T) {
		plugin := &Plugin{
			APIVersion: "v1",
			Info: PluginInfo{
				GatewayMinVersion: "1.0.0",
				Type:              "enterprise",
				Category:          "auth",
				Name:              "test-plugin",
				Image:             "test-image:v1",
				Title:             "Test Plugin",
				TitleI18n:         map[string]string{"zh": "测试插件"},
				Description:       "Test plugin description",
				DescriptionI18n:   map[string]string{"zh": "测试插件描述"},
				IconURL:           "https://example.com/icon.png",
				Version:           "1.0.0",
				Contact: &PluginContact{
					Name:  "Test Author",
					URL:   "https://example.com",
					Email: "test@example.com",
				},
			},
			Spec: PluginSpec{
				Phase:             "AUTHN",
				Priority:          100,
				ConfigSchema:      &ConfigSchema{OpenAPIV3Schema: map[string]interface{}{"type": "object"}},
				RouteConfigSchema: &ConfigSchema{OpenAPIV3Schema: map[string]interface{}{"type": "object"}},
			},
		}

		assert.Equal(t, "v1", plugin.APIVersion)
		assert.Equal(t, "1.0.0", plugin.Info.GatewayMinVersion)
		assert.Equal(t, "enterprise", plugin.Info.Type)
		assert.Equal(t, "auth", plugin.Info.Category)
		assert.Equal(t, "test-plugin", plugin.Info.Name)
		assert.Equal(t, "test-image:v1", plugin.Info.Image)
		assert.Equal(t, "Test Plugin", plugin.Info.Title)
		assert.Equal(t, "测试插件", plugin.Info.TitleI18n["zh"])
		assert.Equal(t, "Test plugin description", plugin.Info.Description)
		assert.Equal(t, "测试插件描述", plugin.Info.DescriptionI18n["zh"])
		assert.Equal(t, "https://example.com/icon.png", plugin.Info.IconURL)
		assert.Equal(t, "1.0.0", plugin.Info.Version)
		assert.NotNil(t, plugin.Info.Contact)
		assert.Equal(t, "Test Author", plugin.Info.Contact.Name)
		assert.Equal(t, "https://example.com", plugin.Info.Contact.URL)
		assert.Equal(t, "test@example.com", plugin.Info.Contact.Email)
		assert.Equal(t, "AUTHN", plugin.Spec.Phase)
		assert.Equal(t, 100, plugin.Spec.Priority)
		assert.NotNil(t, plugin.Spec.ConfigSchema)
		assert.NotNil(t, plugin.Spec.RouteConfigSchema)
	})

	t.Run("PluginContact struct", func(t *testing.T) {
		contact := &PluginContact{
			Name:  "John Doe",
			URL:   "https://johndoe.com",
			Email: "john@example.com",
		}

		assert.Equal(t, "John Doe", contact.Name)
		assert.Equal(t, "https://johndoe.com", contact.URL)
		assert.Equal(t, "john@example.com", contact.Email)
	})

	t.Run("ConfigSchema struct", func(t *testing.T) {
		schema := &ConfigSchema{
			OpenAPIV3Schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"key": map[string]interface{}{"type": "string"},
				},
			},
		}

		assert.NotNil(t, schema.OpenAPIV3Schema)
		assert.Equal(t, "object", schema.OpenAPIV3Schema["type"])
	})
}

func TestPluginInfo_EmptyI18nMaps(t *testing.T) {
	info := &PluginInfo{
		Title:           "Default Title",
		TitleI18n:       map[string]string{},
		Description:     "Default Description",
		DescriptionI18n: map[string]string{},
	}

	t.Run("GetTitle returns default for empty i18n map", func(t *testing.T) {
		result := info.GetTitle("zh")
		assert.Equal(t, "Default Title", result)
	})

	t.Run("GetDescription returns default for empty i18n map", func(t *testing.T) {
		result := info.GetDescription("zh")
		assert.Equal(t, "Default Description", result)
	})
}

func TestPluginSpec_Fields(t *testing.T) {
	t.Run("PluginSpec with all fields", func(t *testing.T) {
		spec := PluginSpec{
			Phase:    "AUTHN",
			Priority: 100,
			ConfigSchema: &ConfigSchema{
				OpenAPIV3Schema: map[string]interface{}{
					"type": "object",
				},
			},
			RouteConfigSchema: &ConfigSchema{
				OpenAPIV3Schema: map[string]interface{}{
					"type": "object",
				},
			},
		}

		assert.Equal(t, "AUTHN", spec.Phase)
		assert.Equal(t, 100, spec.Priority)
		assert.NotNil(t, spec.ConfigSchema)
		assert.NotNil(t, spec.RouteConfigSchema)
	})

	t.Run("PluginSpec with nil schemas", func(t *testing.T) {
		spec := PluginSpec{
			Phase:             "STATS",
			Priority:          50,
			ConfigSchema:      nil,
			RouteConfigSchema: nil,
		}

		assert.Equal(t, "STATS", spec.Phase)
		assert.Equal(t, 50, spec.Priority)
		assert.Nil(t, spec.ConfigSchema)
		assert.Nil(t, spec.RouteConfigSchema)
	})
}

// Helper functions
func boolPtr(b bool) *bool {
	return &b
}

func intPtr(i int) *int {
	return &i
}
