package main

import "fyne.io/fyne/v2"

type ModalWindow struct {
	win            fyne.Window
	enableMenuItem func()
}

func NewModalWindow(win fyne.Window,
	enableMenuItem func()) ModalWindow {
	w := ModalWindow{
		win:            win,
		enableMenuItem: enableMenuItem,
	}
	w.win.SetCloseIntercept(func() {
		w.Hide()
	})
	return w
}

func (w *ModalWindow) Hide() {
	w.enableMenuItem()
	w.win.Hide()
}
