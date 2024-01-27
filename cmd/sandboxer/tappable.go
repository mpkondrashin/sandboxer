package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type tappableIcon struct {
	widget.Icon
	callback func()
}

func newTappableIcon(res fyne.Resource, callback func()) *tappableIcon {
	icon := &tappableIcon{
		callback: callback,
	}
	icon.ExtendBaseWidget(icon)
	icon.SetResource(res)
	return icon
}

func (t *tappableIcon) Tapped(_ *fyne.PointEvent) {
	log.Println("I have been tapped")
	t.callback()
}

func (t *tappableIcon) TappedSecondary(_ *fyne.PointEvent) {
}
