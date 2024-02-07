package main

import (
	"sandboxer/pkg/logging"

	"fyne.io/fyne/v2"
)

type ModalWindow struct {
	win            fyne.Window
	enableMenuItem func()
	quit           func()
}

func NewModalWindow(win fyne.Window,
	enableMenuItem func()) ModalWindow {
	w := ModalWindow{
		win:            win,
		enableMenuItem: enableMenuItem,
	}
	w.win.SetCloseIntercept(w.Hide)
	return w
}

func (w *ModalWindow) SetQuit(quit func()) {
	logging.Debugf("ModalWindow SetQuit: %p", quit)
	w.quit = quit
}

func (w *ModalWindow) Hide() {
	logging.Debugf("ModalWindow Hide, q = %p", w.quit)
	if w.quit != nil {
		w.quit()
	}
	w.enableMenuItem()
	w.win.Hide()
}
