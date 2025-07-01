package ide

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Ide struct {
	app           fyne.App
	window        fyne.Window
	editor        *widget.Entry
	currentFile   string
	currentFolder string
	fileLabel     *widget.Label
	fileTree      *widget.Tree
	treeData      map[string][]string
	splitPanel    *container.Split
}

func NewIde() *Ide {
	aether := app.NewWithID("com.aether.ide")

	aether.SetIcon(nil)

	myWindow := aether.NewWindow("Aether")
	myWindow.Resize(fyne.NewSize(1200, 800))
	myWindow.CenterOnScreen()

	LazyLoadIcons()

	return &Ide{
		app:    aether,
		window: myWindow,
	}
}

func (i *Ide) setupUi() {
	i.editor = widget.NewMultiLineEntry()
	i.editor.Wrapping = fyne.TextWrapWord
	i.editor.SetPlaceHolder("start here")

	i.editor.TextStyle = fyne.TextStyle{
		Monospace: true,
		TabWidth:  4,
	}

	i.fileLabel = widget.NewLabel("untitled.txt")
	i.fileLabel.TextStyle = fyne.TextStyle{
		Bold: true,
	}

	editorContainer := container.NewBorder(
		nil, nil, nil, nil,
		container.NewScroll(i.editor),
	)

	i.treeData = make(map[string][]string)
	i.treeData[""] = []string{"root"}
	i.treeData["root"] = []string{}
	i.fileTree = i.createFileTree()

	i.fileTree.Resize(fyne.NewSize(240, 400))

	newFileBtn := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), i.newFile)
	newFileBtn.SetText("")
	newFolderBtn := widget.NewButtonWithIcon("", theme.FolderIcon(), i.newFolder)
	newFolderBtn.SetText("")
	openFolderBtn := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), i.openFolder)
	openFolderBtn.SetText("")

	toolbar := container.NewHBox(newFileBtn, newFolderBtn, openFolderBtn)

	treeContainer := container.NewBorder(
		toolbar, nil, nil, nil,
		container.NewScroll(i.fileTree),
	)

	treeContainer.Resize(fyne.NewSize(250, 600))

	editorArea := container.NewBorder(
		nil,
		i.fileLabel,
		nil,
		nil,
		editorContainer,
	)

	i.splitPanel = container.NewHSplit(treeContainer, editorArea)
	i.splitPanel.Offset = 0.2
	i.splitPanel.SetOffset(0.2)

	i.window.SetContent(i.splitPanel)
	i.setupEventHandlers()
	i.setupMenu()
}

func (i *Ide) setupEventHandlers() {
	newFileItem := fyne.NewMenuItem("New File ...", i.newFile)
	newFileItem.Icon = theme.DocumentCreateIcon()

	newFolderItem := fyne.NewMenuItem("New Folder", i.newFolder)
	newFolderItem.Icon = theme.FolderIcon()

	openFolderItem := fyne.NewMenuItem("Open Folder ...", i.openFolder)
	openFolderItem.Icon = theme.FolderOpenIcon()

	openItem := fyne.NewMenuItem("Open File ...", i.openFile)
	openItem.Icon = theme.DocumentIcon()

	saveItem := fyne.NewMenuItem("Save File", i.saveFile)
	saveItem.Icon = theme.DocumentSaveIcon()

	saveItemAs := fyne.NewMenuItem("Save As ...", i.saveFileAs)

	quitItem := fyne.NewMenuItem("Quit", func() {
		i.app.Quit()
	})

	fileMenu := fyne.NewMenu("File", newFileItem, newFolderItem, fyne.NewMenuItemSeparator(),
		openFolderItem, openItem, fyne.NewMenuItemSeparator(),
		saveItem, saveItemAs, fyne.NewMenuItemSeparator(), quitItem)

	mainMenu := fyne.NewMainMenu(fileMenu)
	i.window.SetMainMenu(mainMenu)
}

func (i *Ide) setupMenu() {
	i.window.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		if key.Name == fyne.KeyS && (key.Physical.ScanCode == 0 || key.Physical.ScanCode == 1) {
			i.saveFile()
		}
	})
}

func (i *Ide) newFile() {
	if i.currentFolder == "" {
		dialog.ShowInformation("No Folder", "Please open a folder first to create a new file", i.window)
		return
	}

	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter filename...")

	dialog.ShowForm("New File", "Create", "Cancel", []*widget.FormItem{
		widget.NewFormItem("Filename", entry),
	}, func(confirm bool) {
		if !confirm || entry.Text == "" {
			return
		}

		filePath := filepath.Join(i.currentFolder, entry.Text)

		if _, err := os.Stat(filePath); err == nil {
			dialog.ShowError(fmt.Errorf("file %s already exists", entry.Text), i.window)
			return
		}

		err := os.WriteFile(filePath, []byte(""), 0644)
		if err != nil {
			dialog.ShowError(err, i.window)
			return
		}

		i.loadFolderContents()

		i.openFileByPath(filePath)
	}, i.window)
}

func (i *Ide) newFolder() {
	if i.currentFolder == "" {
		dialog.ShowInformation("No Folder", "Please open a folder first to create a new folder", i.window)
		return
	}

	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter folder name...")

	dialog.ShowForm("New Folder", "Create", "Cancel", []*widget.FormItem{
		widget.NewFormItem("Folder Name", entry),
	}, func(confirm bool) {
		if !confirm || entry.Text == "" {
			return
		}

		folderPath := filepath.Join(i.currentFolder, entry.Text)

		if _, err := os.Stat(folderPath); err == nil {
			dialog.ShowError(fmt.Errorf("folder %s already exists", entry.Text), i.window)
			return
		}

		err := os.MkdirAll(folderPath, 0755)
		if err != nil {
			dialog.ShowError(err, i.window)
			return
		}

		i.loadFolderContents()
	}, i.window)
}

func (i *Ide) openFile() {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, i.window)
			return
		}
		if reader == nil {
			return
		}
		defer reader.Close()

		data, err := io.ReadAll(reader)
		if err != nil {
			dialog.ShowError(err, i.window)
			return
		}

		i.editor.SetText(string(data))
		i.currentFile = reader.URI().Path()
		i.fileLabel.SetText(filepath.Base(i.currentFile))
		i.window.SetTitle(fmt.Sprintf("Aether - %s", filepath.Base(i.currentFile)))
	}, i.window)
}

func (i *Ide) saveFile() {
	if i.currentFile == "" {
		i.saveFileAs()
		return
	}

	err := os.WriteFile(i.currentFile, []byte(i.editor.Text), 0644)
	if err != nil {
		dialog.ShowError(err, i.window)
		return
	}

	dialog.ShowInformation("Saved", fmt.Sprintf("File saved to %s", i.currentFile), i.window)
}

func (i *Ide) saveFileAs() {
	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, i.window)
			return
		}

		if writer == nil {
			return
		}
		defer writer.Close()

		_, err = writer.Write([]byte(i.editor.Text))
		if err != nil {
			dialog.ShowError(err, i.window)
			return
		}

		i.currentFile = writer.URI().Path()
		i.fileLabel.SetText(filepath.Base(i.currentFile))
		i.window.SetTitle(fmt.Sprintf("Aether - %s", filepath.Base(i.currentFile)))

		dialog.ShowInformation("Saved", fmt.Sprintf("File saved to %s", i.currentFile), i.window)
	}, i.window)
}

func (i *Ide) createFileTree() *widget.Tree {
	tree := widget.NewTree(
		func(uid string) []string {
			children := i.treeData[uid]
			return children
		},
		func(uid string) bool {
			children, ok := i.treeData[uid]
			return ok && len(children) > 0
		},
		func(branch bool) fyne.CanvasObject {
			icon := widget.NewIcon(theme.DocumentIcon())
			icon.Resize(fyne.NewSize(16, 16))
			label := widget.NewLabel("Template Object")
			return container.NewHBox(icon, label)
		},
		func(uid string, branch bool, obj fyne.CanvasObject) {
			hbox := obj.(*fyne.Container)
			objects := hbox.Objects
			if len(objects) >= 2 {
				icon := objects[0].(*widget.Icon)
				label := objects[1].(*widget.Label)

				if uid == "" {
					return
				}

				if uid == "root" {
					label.SetText(i.getFolderDisplayName())
				} else {
					label.SetText(filepath.Base(uid))
				}

				if branch {
					icon.SetResource(theme.FolderIcon())
				} else {
					icon.SetResource(i.getFileIcon(uid))
				}
			}
		},
	)

	tree.OnSelected = func(uid string) {
		if uid == "" || uid == "root" {
			return
		}

		fullPath := uid
		if !filepath.IsAbs(uid) {
			fullPath = filepath.Join(i.currentFolder, uid)
		}

		if info, err := os.Stat(fullPath); err == nil {
			if info.IsDir() {
				fyne.Do(func() {
					if tree.IsBranchOpen(uid) {
						tree.CloseBranch(uid)
					} else {
						tree.OpenBranch(uid)
					}
					tree.UnselectAll()
				})
			} else {
				i.openFileByPath(fullPath)
			}
		}
	}

	return tree
}

func (i *Ide) getFolderDisplayName() string {
	if i.currentFolder == "" {
		return "Open a folder to get started"
	}
	return filepath.Base(i.currentFolder)
}

func (i *Ide) getFileIcon(uid string) fyne.Resource {
	return GetFileIcon(uid)
}

func (i *Ide) loadFolderContents() {
	if i.currentFolder == "" {
		i.treeData = make(map[string][]string)
		i.treeData[""] = []string{"root"}
		i.treeData["root"] = []string{}
		i.fileTree.Refresh()
		return
	}

	// fmt.Printf("Loading folder contents for: %s\n", i.currentFolder)
	i.treeData = make(map[string][]string)

	i.treeData[""] = []string{"root"}
	rootChildren := []string{}

	entries, err := os.ReadDir(i.currentFolder)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		dialog.ShowError(err, i.window)
		return
	}

	// fmt.Printf("Found %d entries\n", len(entries))

	var dirs, files []os.DirEntry
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry)
		} else {
			files = append(files, entry)
		}
	}

	for _, dir := range dirs {
		dirPath := filepath.Join(i.currentFolder, dir.Name())
		rootChildren = append(rootChildren, dirPath)
		i.loadSubdirectory(dirPath)
	}

	for _, file := range files {
		filePath := filepath.Join(i.currentFolder, file.Name())
		rootChildren = append(rootChildren, filePath)
	}

	// fmt.Printf("Root children count: %d\n", len(rootChildren))
	// for i, child := range rootChildren {
	// 	fmt.Printf("Child %d: %s\n", i, child)
	// }

	i.treeData["root"] = rootChildren
	i.fileTree.Refresh()

	if len(rootChildren) > 0 {
		i.fileTree.OpenBranch("root")
	}
}

func (i *Ide) loadSubdirectory(dirPath string) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return
	}

	children := []string{}
	for _, entry := range entries {
		childPath := filepath.Join(dirPath, entry.Name())
		children = append(children, childPath)

		if entry.IsDir() {
			i.loadSubdirectory(childPath)
		}
	}

	i.treeData[dirPath] = children
}

func (i *Ide) openFolder() {
	dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
		if err != nil {
			dialog.ShowError(err, i.window)
			return
		}
		if uri == nil {
			return
		}

		i.currentFolder = uri.Path()
		// fmt.Printf("Opening folder: %s\n", i.currentFolder)
		i.loadFolderContents()
		i.window.SetTitle(fmt.Sprintf("Aether - %s", filepath.Base(i.currentFolder)))
	}, i.window)
}

func (i *Ide) openFileByPath(filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		dialog.ShowError(err, i.window)
		return
	}

	i.editor.SetText(string(data))
	i.currentFile = filePath
	i.fileLabel.SetText(filepath.Base(i.currentFile))
	i.window.SetTitle(fmt.Sprintf("Aether - %s", filepath.Base(i.currentFile)))
}

func (i *Ide) Run() {
	i.setupUi()

	i.window.ShowAndRun()
}
