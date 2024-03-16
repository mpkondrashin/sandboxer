/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

context_menu.go

Icon with context menu
*/
package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type contextMenuIcon struct {
	widget.Icon
	Menu *fyne.Menu
}

func (b *contextMenuIcon) Tapped(e *fyne.PointEvent) {
	//logging.Debugf("%v contextMenuIcon Tapped(%v)", b, e)
	widget.ShowPopUpMenuAtPosition(b.Menu, fyne.CurrentApp().Driver().CanvasForObject(b), e.AbsolutePosition)
}

func newContextMenuIcon(res fyne.Resource, menu *fyne.Menu) *contextMenuIcon {
	b := &contextMenuIcon{
		Menu: menu,
	}
	b.Resource = res
	b.ExtendBaseWidget(b)
	return b
}
