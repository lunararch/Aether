# Aether IDE

A modern Go-based IDE built with Fyne, featuring a professional file tree sidebar with custom icons.

## Features

### File Tree Sidebar
- **Expandable/collapsible folder navigation** - Click on folders to expand/collapse them
- **Quick action buttons** - New file, new folder, and open folder buttons
- **Custom file type icons** - Professional icons for different file types and extensions
- **Resizable and scrollable** - Sidebar can be resized and scrolls when content overflows
- **Click-to-expand** - Click anywhere on a folder to expand/collapse it (not just the arrow)

### Supported File Types & Icons
The IDE supports custom icons for numerous file types including:

**Programming Languages:**
- TypeScript (.ts) - typescript icons
- JavaScript (.js) - javascript icons
- React (.jsx, .tsx) - react icons
- Vue.js (.vue) - vuejs icons
- Go (.go) - go icons
- Python (.py) - python icons
- Java (.java) - java icons
- C (.c) - c icons
- C++ (.cpp, .cc, .cxx, .c++) - cplusplus icons
- C# (.cs) - csharp icons
- PHP (.php) - php icons
- Ruby (.rb) - ruby icons
- Rust (.rs) - rust icons
- Swift (.swift) - swift icons
- Kotlin (.kt, .kts) - kotlin icons
- Dart (.dart) - dart icons

**Web Technologies:**
- HTML (.html, .htm) - html5 icons
- CSS (.css) - css3 icons
- Sass (.scss, .sass) - sass icons
- Less (.less) - less icons

**Data & Config:**
- JSON (.json) - json icons
- XML (.xml) - xml icons
- YAML (.yaml, .yml) - yaml icons
- Markdown (.md) - markdown icons

**Database:**
- SQL (.sql) - mysql icons

**Scripts & DevOps:**
- Bash (.sh, .bash) - bash icons
- PowerShell (.ps1) - powershell icons
- Docker (dockerfile, .dockerfile) - docker icons

**Special Files:**
- package.json, package-lock.json - nodejs icons
- go.mod, go.sum - go icons
- cargo.toml, cargo.lock - rust icons
- requirements.txt, pyproject.toml - python icons
- tsconfig.json - typescript icons
- And many more...

## Architecture

### File Structure
```
main.go              # Main application entry point
ide/
  ide.go            # Main IDE logic, UI, and file tree implementation
  icons.go          # Icon detection and loading logic
icons/              # Large directory of custom icons (SVG format)
  typescript/
  javascript/
  go/
  python/
  ...
```

### Key Components

**IDE Core (`ide/ide.go`)**
- Main IDE window and layout
- File tree implementation using Fyne's Tree widget
- Event handlers for file/folder interactions
- UI threading management with `fyne.Do()`

**Icon System (`ide/icons.go`)**
- `loadIconFromFile()` - Loads SVG icons from the icons directory
- `getFileIcon()` - Determines appropriate icon based on file extension/name
- Automatic fallback to theme icons when custom icons aren't available
- Support for multiple icon variants (original, plain, line)

## Technical Details

### Icon Loading
The icon system tries multiple variants for each file type in order of preference:
1. `{iconName}-original.svg`
2. `{iconName}-plain.svg`
3. `{iconName}-line.svg`
4. `{iconName}.svg`

If no custom icon is found, it falls back to Fyne's built-in theme icons.

### UI Threading
All UI operations are properly threaded using `fyne.Do()` to avoid Fyne UI thread errors.

### Performance
- Icons are loaded on-demand
- File tree data is cached to avoid unnecessary file system operations
- Efficient SVG resource handling using `fyne.NewStaticResource()`

## Recent Fixes

### Resolved Issues
- ✅ Fixed SVG icons appearing as white squares
- ✅ Resolved Fyne UI thread errors
- ✅ Improved icon loading reliability
- ✅ Cleaned up code organization and removed unused imports
- ✅ Enhanced file type detection and icon mapping

### Changes Made
- Switched from `fyne.LoadResourceFromPath()` to `fyne.NewStaticResource()` for better SVG handling
- Refactored icon logic into separate `icons.go` file for modularity
- Added comprehensive file type and extension mapping
- Implemented proper error handling and fallbacks
- Added icon size constraints for consistent display

## Usage

To run the IDE:
```bash
go run .
```

The IDE will open with a file tree sidebar showing the current directory. You can:
- Click on folders to expand/collapse them
- Use the quick action buttons to create new files/folders or open different directories
- Resize the sidebar by dragging the splitter
- Scroll through large directory structures

## Dependencies

- [Fyne v2](https://fyne.io/) - Cross-platform GUI toolkit for Go
- Standard Go libraries (os, path/filepath, strings)

The IDE uses a large collection of developer icons (devicons) in SVG format for professional file type recognition.
