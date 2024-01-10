package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type PageFolder struct {
	folderEntry *widget.Entry
}

var _ Page = &PageFolder{}

func (p *PageFolder) Name() string {
	return "Folder"
}

func (p *PageFolder) Content(win fyne.Window, model *Model) fyne.CanvasObject {
	labelFolder := widget.NewLabel("Installation folder:")
	p.folderEntry = widget.NewEntry()
	p.folderEntry.SetText(model.config.Folder)
	folderButton := widget.NewButton("Choose...", func() {
		folderDialog := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri == nil {
				return
			}
			p.folderEntry.SetText(uri.Path())
		}, win)
		folderDialog.Show()
	})
	return container.NewVBox(labelFolder, p.folderEntry, folderButton)
}

func (p *PageFolder) AquireData(model *Model) error {
	model.config.Folder = p.folderEntry.Text
	return nil
}
