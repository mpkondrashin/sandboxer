package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type contextMenuIcon struct {
	widget.Icon
	menu *fyne.Menu
}

func (b *contextMenuIcon) Tapped(e *fyne.PointEvent) {
	widget.ShowPopUpMenuAtPosition(b.menu, fyne.CurrentApp().Driver().CanvasForObject(b), e.AbsolutePosition)
}

func newContextMenuIcon(res fyne.Resource, menu *fyne.Menu) *contextMenuIcon {
	b := &contextMenuIcon{
		menu: menu,
	}
	b.Resource = res
	b.ExtendBaseWidget(b)
	return b
}
