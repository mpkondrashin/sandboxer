/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
This software is distributed under MIT license as stated in LICENSE file

page_folder.go

Pick destination folder
*/
package main

import (
	"sandboxer/pkg/globals"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type PageFolder struct {
	BasePage
	folderEntry *widget.Entry
}

var _ Page = &PageFolder{}

func (p *PageFolder) Name() string {
	return "Destination"
}

func (p *PageFolder) Next(previousPage PageIndex) PageIndex {
	p.SavePrevious(previousPage)
	return pgInstallation
}

func (p *PageFolder) Content() fyne.CanvasObject {
	labelFolder := widget.NewLabel("Base folder to install " + globals.AppName + ":")
	p.folderEntry = widget.NewEntry()
	p.folderEntry.SetText(p.wiz.installer.config.Folder)
	folderButton := widget.NewButton("Change...", func() {
		folderDialog := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri == nil {
				return
			}
			p.folderEntry.SetText(uri.Path())
		}, p.wiz.win)
		folderDialog.Show()
	})
	return container.NewVBox(labelFolder,
		container.NewBorder(nil, nil, nil, folderButton, p.folderEntry)) // p.folderEntry, folderButton)
}

func (p *PageFolder) Run() {}

func (p *PageFolder) AquireData(installer *Installer) error {
	installer.config.Folder = p.folderEntry.Text
	return nil
}
