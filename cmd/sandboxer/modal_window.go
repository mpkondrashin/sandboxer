/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

modal_window.go

Window that automatically disables and enables menu item
*/
package main

import (
	"fyne.io/fyne/v2"
)

type ModalWindowContent interface {
	Name() string
	Icon() fyne.Resource
	Content(w *ModalWindow) fyne.CanvasObject
	Show()
	Hide()
}

type ModalWindow struct {
	win      fyne.Window
	trayApp  *TrayApp
	MenuItem *fyne.MenuItem
	content  ModalWindowContent
}

func NewModalWindow(content ModalWindowContent, trayApp *TrayApp) *ModalWindow {
	w := &ModalWindow{
		content: content,
		win:     trayApp.app.NewWindow(content.Name()),
		trayApp: trayApp,
	}

	w.MenuItem = fyne.NewMenuItem(content.Name()+"...", w.Show)
	w.MenuItem.Icon = content.Icon()
	w.win.SetCloseIntercept(w.Hide)
	w.win.SetContent(content.Content(w))
	return w
}

func (w *ModalWindow) Show() {
	w.content.Show()
	w.MenuItem.Disabled = true
	w.win.Show()
	w.trayApp.menu.Refresh()
}

func (w *ModalWindow) Hide() {
	w.content.Hide()
	w.MenuItem.Disabled = false
	w.win.Hide()
	w.trayApp.menu.Refresh()
}
