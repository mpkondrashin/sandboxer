package fatal

import (
	"sandboxer/pkg/globals"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func Warning(title, text string) {
	//dialog.ShowInformation("Information", "This is a sample message", s.window)
	a := app.New()
	w := a.NewWindow(globals.AppName + " " + title)
	//w.SetMaster() // will it exit the application?
	w.SetContent(container.NewVBox(
		widget.NewLabel(text),
		widget.NewButton("Ok", func() { a.Quit() }),
	))
	w.ShowAndRun()
}
