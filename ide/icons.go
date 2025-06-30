package ide

import (
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// loadIconFromFile loads an icon from the icons folder
func (i *Ide) loadIconFromFile(iconName string) fyne.Resource {
	// Try different icon variants in order of preference
	variants := []string{
		iconName + "-original-wordmark.svg",
		iconName + "-original.svg",
		iconName + "-plain.svg",
		iconName + "-line.svg",
		iconName + ".svg",
	}

	for _, variant := range variants {
		iconPath := filepath.Join("icons", iconName, variant)
		if _, err := os.Stat(iconPath); err == nil {
			// Try reading the file and creating a resource
			data, err := os.ReadFile(iconPath)
			if err == nil {
				resource := fyne.NewStaticResource(variant, data)
				return resource
			} else {
				// Could not read file
			}
		}
	}

	// Fallback to theme icon
	return theme.DocumentIcon()
}

// getFileIcon returns appropriate icon based on file extension and filename
func (i *Ide) getFileIcon(filePath string) fyne.Resource {
	ext := strings.ToLower(filepath.Ext(filePath))
	filename := strings.ToLower(filepath.Base(filePath))

	// Handle special filenames first
	switch filename {
	case "dockerfile", "dockerfile.dev", "dockerfile.prod":
		return i.loadIconFromFile("docker")
	case ".gitignore", ".gitattributes", ".gitmodules":
		return i.loadIconFromFile("git")
	case "package.json", "package-lock.json":
		return i.loadIconFromFile("nodejs")
	case "cargo.toml", "cargo.lock":
		return i.loadIconFromFile("rust")
	case "go.mod", "go.sum":
		return i.loadIconFromFile("go")
	case "requirements.txt", "pyproject.toml":
		return i.loadIconFromFile("python")
	case "pom.xml":
		return i.loadIconFromFile("java")
	case "composer.json", "composer.lock":
		return i.loadIconFromFile("php")
	case "yarn.lock":
		return i.loadIconFromFile("yarn")
	case "readme.md", "readme":
		return i.loadIconFromFile("markdown")
	case "makefile":
		return i.loadIconFromFile("cmake")
	case "webpack.config.js", "webpack.config.ts":
		return i.loadIconFromFile("webpack")
	case "vite.config.js", "vite.config.ts":
		return i.loadIconFromFile("vite")
	case "tsconfig.json":
		return i.loadIconFromFile("typescript")
	}

	// Handle file extensions
	switch ext {
	case ".ts":
		return i.loadIconFromFile("typescript")
	case ".js":
		return i.loadIconFromFile("javascript")
	case ".jsx":
		return i.loadIconFromFile("react")
	case ".tsx":
		return i.loadIconFromFile("react")
	case ".vue":
		return i.loadIconFromFile("vuejs")
	case ".go":
		return i.loadIconFromFile("go")
	case ".py":
		return i.loadIconFromFile("python")
	case ".java":
		return i.loadIconFromFile("java")
	case ".c":
		return i.loadIconFromFile("c")
	case ".cpp", ".cc", ".cxx", ".c++":
		return i.loadIconFromFile("cplusplus")
	case ".cs":
		return i.loadIconFromFile("csharp")
	case ".php":
		return i.loadIconFromFile("php")
	case ".rb":
		return i.loadIconFromFile("ruby")
	case ".rs":
		return i.loadIconFromFile("rust")
	case ".swift":
		return i.loadIconFromFile("swift")
	case ".kt", ".kts":
		return i.loadIconFromFile("kotlin")
	case ".dart":
		return i.loadIconFromFile("dart")
	case ".html", ".htm":
		return i.loadIconFromFile("html5")
	case ".css":
		return i.loadIconFromFile("css3")
	case ".scss", ".sass":
		return i.loadIconFromFile("sass")
	case ".less":
		return i.loadIconFromFile("less")
	case ".json":
		return i.loadIconFromFile("json")
	case ".xml":
		return i.loadIconFromFile("xml")
	case ".yaml", ".yml":
		return i.loadIconFromFile("yaml")
	case ".md":
		return i.loadIconFromFile("markdown")
	case ".sql":
		return i.loadIconFromFile("mysql")
	case ".sh", ".bash":
		return i.loadIconFromFile("bash")
	case ".ps1":
		return i.loadIconFromFile("powershell")
	case ".dockerfile":
		return i.loadIconFromFile("docker")
	default:
		return theme.DocumentIcon()
	}
}
