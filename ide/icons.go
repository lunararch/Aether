package ide

import (
	"embed"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

//go:embed icons/*.svg
var iconFS embed.FS

// FileIconMap maps file extensions to their corresponding SVG icon resource
var FileIconMap map[string]fyne.Resource

func init() {
	FileIconMap = make(map[string]fyne.Resource)

	// Initialize file extension to icon mappings
	FileIconMap[".go"] = loadIcon("go.svg")
	FileIconMap[".py"] = loadIcon("python.svg")
	FileIconMap[".js"] = loadIcon("javascript.svg")
	FileIconMap[".ts"] = loadIcon("typescript.svg")
	FileIconMap[".jsx"] = loadIcon("react.svg")
	FileIconMap[".tsx"] = loadIcon("react.svg")
	FileIconMap[".html"] = loadIcon("html5.svg")
	FileIconMap[".css"] = loadIcon("css3.svg")
	FileIconMap[".json"] = loadIcon("json.svg")
	FileIconMap[".xml"] = loadIcon("xml.svg")
	FileIconMap[".yaml"] = loadIcon("yaml.svg")
	FileIconMap[".yml"] = loadIcon("yaml.svg")
	FileIconMap[".md"] = loadIcon("markdown.svg")
	FileIconMap[".java"] = loadIcon("java.svg")
	FileIconMap[".c"] = loadIcon("c.svg")
	FileIconMap[".cpp"] = loadIcon("cplusplus.svg")
	FileIconMap[".cs"] = loadIcon("csharp.svg")
	FileIconMap[".sh"] = loadIcon("bash.svg")
	FileIconMap[".rb"] = loadIcon("ruby.svg")
	FileIconMap[".rs"] = loadIcon("rust.svg")
	FileIconMap[".kt"] = loadIcon("kotlin.svg")
	FileIconMap[".lua"] = loadIcon("lua.svg")
	FileIconMap[".dockerfile"] = loadIcon("docker.svg")
	FileIconMap[".svelte"] = loadIcon("svelte.svg")
	FileIconMap[".vue"] = loadIcon("vuejs.svg")
	FileIconMap[".dart"] = loadIcon("flutter.svg")
	FileIconMap[".gradle"] = loadIcon("gradle.svg")
	FileIconMap[".gitignore"] = loadIcon("git.svg")
}

// loadIcon loads an SVG icon from the embedded filesystem
func loadIcon(name string) fyne.Resource {
	data, err := iconFS.ReadFile("icons/" + name)
	if err != nil {
		return theme.FileIcon()
	}

	resource := fyne.NewStaticResource(name, data)
	return resource
}

// GetFileIcon returns the appropriate icon for a file based on its extension
func GetFileIcon(filePath string) fyne.Resource {
	ext := strings.ToLower(filepath.Ext(filePath))

	// Special case for Dockerfile which doesn't have an extension
	if strings.ToLower(filepath.Base(filePath)) == "dockerfile" {
		return FileIconMap[".dockerfile"]
	}

	// Special case for .gitignore
	if strings.ToLower(filepath.Base(filePath)) == ".gitignore" {
		return FileIconMap[".gitignore"]
	}

	// Check if we have a specific icon for this extension
	if icon, ok := FileIconMap[ext]; ok {
		return icon
	}

	// Default to standard file icon
	return theme.FileIcon()
}
