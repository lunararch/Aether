package ide

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"

	//"fyne.io/fyne/v2/theme"
	"io"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2/widget"
)

type Ide struct {
	app           fyne.App
	window        fyne.Window
	editor        *widget.Entry
	currentFile   string
	currentFolder string
	fileLabel     *widget.Label
	splitPanel    *container.Split
}

func NewIde() *Ide {
	aether := app.NewWithID("com.aether.ide")

	aether.SetIcon(nil)

	myWindow := aether.NewWindow("Aether")
	myWindow.Resize(fyne.NewSize(1200, 800))
	myWindow.CenterOnScreen()

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

	content := container.NewBorder(
		nil,
		i.fileLabel,
		nil,
		nil,
		editorContainer,
	)

	i.window.SetContent(content)
	i.setupEventHandlers()
	i.setupMenu()
}

func (i *Ide) setupEventHandlers() {
	newFileItem := fyne.NewMenuItem("New File ...", i.newFile)
	//newFileItem.Icon = theme.DocumentCreateIcon()

	newFolderItem := fyne.NewMenuItem("New Folder", i.newFolder)

	openItem := fyne.NewMenuItem("Open ...", i.openFile)
	//openItem.Icon = theme.FolderOpenIcon()

	saveItem := fyne.NewMenuItem("Save File", i.saveFile)
	//saveItem.Icon = theme.DocumentSaveIcon()

	saveItemAs := fyne.NewMenuItem("Save As ...", i.saveFileAs)

	quitItem := fyne.NewMenuItem("Quit", func() {
		i.app.Quit()
	})

	fileMenu := fyne.NewMenu("File", newFileItem, newFolderItem, openItem, fyne.NewMenuItemSeparator(),
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
	i.editor.SetText("")
	i.currentFile = ""
	i.fileLabel.SetText("unititled.txt")
	i.window.SetTitle("Aether")
}

func (i *Ide) newFolder() {

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

func (i *Ide) Run() {
	i.setupUi()

	i.window.ShowAndRun()
}
