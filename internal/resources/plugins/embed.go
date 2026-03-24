// Package plugins provides embedded plugin resources for the SDK.
package plugins

import (
	"embed"
	"io/fs"
	"strings"
)

//go:embed plugins.properties
var pluginsProperties []byte

//go:embed *
var pluginsFS embed.FS

// GetPluginsProperties returns the plugins.properties content.
func GetPluginsProperties() []byte {
	return pluginsProperties
}

// GetPluginsPropertiesString returns the plugins.properties content as string.
func GetPluginsPropertiesString() string {
	return string(pluginsProperties)
}

// GetPluginSpec returns the spec.yaml content for a plugin.
func GetPluginSpec(pluginName string) ([]byte, error) {
	path := pluginName + "/spec.yaml"
	return pluginsFS.ReadFile(path)
}

// GetPluginReadme returns the README content for a plugin.
// It tries README.md, README_CN.md, README_EN.md in order.
func GetPluginReadme(pluginName string, lang string) ([]byte, error) {
	// Try language-specific readme first
	if lang != "" {
		langFile := pluginName + "/README_" + lang + ".md"
		if content, err := pluginsFS.ReadFile(langFile); err == nil {
			return content, nil
		}
	}

	// Try default README.md
	path := pluginName + "/README.md"
	if content, err := pluginsFS.ReadFile(path); err == nil {
		return content, nil
	}

	// Try README_CN.md as fallback for Chinese
	if strings.HasPrefix(lang, "zh") {
		path := pluginName + "/README_CN.md"
		if content, err := pluginsFS.ReadFile(path); err == nil {
			return content, nil
		}
	}

	// Try README_EN.md as fallback for English
	if strings.HasPrefix(lang, "en") {
		path := pluginName + "/README_EN.md"
		if content, err := pluginsFS.ReadFile(path); err == nil {
			return content, nil
		}
	}

	return nil, fs.ErrNotExist
}

// GetPluginIcon returns the icon.png content for a plugin.
func GetPluginIcon(pluginName string) ([]byte, error) {
	path := pluginName + "/icon.png"
	return pluginsFS.ReadFile(path)
}

// ListPlugins returns a list of all plugin names from plugins.properties.
func ListPlugins() []string {
	content := string(pluginsProperties)
	lines := strings.Split(content, "\n")
	var plugins []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			plugins = append(plugins, strings.TrimSpace(parts[0]))
		}
	}
	return plugins
}

// GetPluginImageURL returns the image URL for a plugin from plugins.properties.
func GetPluginImageURL(pluginName string) string {
	content := string(pluginsProperties)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			name := strings.TrimSpace(parts[0])
			if name == pluginName {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return ""
}
