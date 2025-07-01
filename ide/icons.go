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

var FileIconMap map[string]fyne.Resource

func LazyLoadIcons() {
	if FileIconMap != nil {
		return
	}

	FileIconMap = make(map[string]fyne.Resource)

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

func loadIcon(name string) fyne.Resource {
	data, err := iconFS.ReadFile("icons/" + name)
	if err != nil {
		return theme.FileIcon()
	}

	resource := fyne.NewStaticResource(name, data)
	return resource
}

func GetFileIcon(filePath string) fyne.Resource {
	LazyLoadIcons()

	ext := strings.ToLower(filepath.Ext(filePath))
	filename := filepath.Base(filePath)

	switch strings.ToLower(filename) {
	case "dockerfile":
		return FileIconMap[".dockerfile"]
	case ".gitignore":
		return FileIconMap[".gitignore"]
	}

	if strings.ToLower(filename) == ".gitignore" {
		return FileIconMap[".gitignore"]
	}

	if icon, ok := FileIconMap[ext]; ok {
		return icon
	}

	return theme.FileIcon()
}
