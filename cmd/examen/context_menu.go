package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

/*
Does not work:

	type contextMenuContainer struct {
		*fyne.Container
		menu *fyne.Menu
	}

	func (b *contextMenuContainer) Tapped(e *fyne.PointEvent) {
		widget.ShowPopUpMenuAtPosition(b.menu, fyne.CurrentApp().Driver().CanvasForObject(b), e.AbsolutePosition)
	}

	func newContextMenuContainer(object fyne.CanvasObject, menu *fyne.Menu) *contextMenuContainer {
		b := &contextMenuContainer{
			Container: container.NewWithoutLayout(object),
			menu:      menu}
		b.Add(object)
		return b
	}
*/
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

type contextMenuLabel struct {
	widget.Label
	menu *fyne.Menu
}

func (b *contextMenuLabel) Tapped(e *fyne.PointEvent) {
	widget.ShowPopUpMenuAtPosition(b.menu, fyne.CurrentApp().Driver().CanvasForObject(b), e.AbsolutePosition)
}

func newContextMenuLable(label string, menu *fyne.Menu) *contextMenuLabel {
	b := &contextMenuLabel{menu: menu}
	b.Text = label
	b.ExtendBaseWidget(b)
	return b
}
