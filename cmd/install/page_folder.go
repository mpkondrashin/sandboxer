/*
TunnelEffect (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
This software is distributed under MIT license as stated in LICENSE file

page_folder.go

Pick destination folder
*/
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
	return "Destination"
}

func (p *PageFolder) Content(win fyne.Window, installer *Installer) fyne.CanvasObject {
	labelFolder := widget.NewLabel("Installation folder:")
	p.folderEntry = widget.NewEntry()
	p.folderEntry.SetText(installer.config.Folder)
	folderButton := widget.NewButton("Change...", func() {
		folderDialog := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri == nil {
				return
			}
			p.folderEntry.SetText(uri.Path())
		}, win)
		folderDialog.Show()
	})
	return container.NewVBox(labelFolder,
		container.NewBorder(nil, nil, nil, folderButton, p.folderEntry)) // p.folderEntry, folderButton)
}

func (p *PageFolder) Run(win fyne.Window, installer *Installer) {

}

func (p *PageFolder) AquireData(installer *Installer) error {
	installer.config.Folder = p.folderEntry.Text
	return nil
}
