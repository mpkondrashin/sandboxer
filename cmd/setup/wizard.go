package main

import "fyne.io/fyne/v2"

type Page interface {
	Name() string
	Content(win fyne.Window, model *Model) fyne.CanvasObject
	AquireData(model *Model) error
}
