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

type ModalWindow struct {
	win            fyne.Window
	enableMenuItem func()
	quit           func()
}

func NewModalWindow(win fyne.Window, enableMenuItem func()) ModalWindow {
	w := ModalWindow{
		win:            win,
		enableMenuItem: enableMenuItem,
	}
	w.win.SetCloseIntercept(w.Hide)
	return w
}

func (w *ModalWindow) SetQuit(quit func()) {
	w.quit = quit
}

func (w *ModalWindow) Hide() {
	if w.quit != nil {
		w.quit()
	}
	w.enableMenuItem()
	w.win.Hide()
}
