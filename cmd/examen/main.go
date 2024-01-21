package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"

	"examen/pkg/grpc"
	"examen/pkg/state"
	"examen/pkg/task"
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

type Status interface {
	Get(from int, count int)
}

type StatusWindow struct {
	win   fyne.Window
	vbox  *fyne.Container
	from  int32
	count int32
	icons []fyne.Resource
}

func NewStatusWindow(app fyne.App) *StatusWindow {
	s := &StatusWindow{
		win:   app.NewWindow("Status"),
		from:  0,
		count: 10,
	}
	for st := state.StateUnknown; st < state.StateCount; st++ {
		r, err := fyne.LoadResourceFromPath(IconPath(st))
		if err != nil {
			panic(err)
		}
		s.icons = append(s.icons, r)
	}
	s.vbox = container.NewVBox()
	s.win.SetCloseIntercept(func() {
		s.win.Hide()
	})
	return s
}

func (s *StatusWindow) Update() {
	s.vbox.RemoveAll()
	err := grpc.Status(s.from, s.count, func(tsk *task.Task) {
		pathLabel := widget.NewLabel(tsk.Path)
		//var icon *tappableIcon
		icon := newTappableIcon(s.icons[tsk.State], func() {
			// pup up menu
		})
		line := container.NewBorder(nil, nil, container.NewHBox(icon, pathLabel), nil)
		s.vbox.Add(line)
	})
	if err != nil {
		s.win.SetContent(widget.NewLabel(err.Error()))
		//dialog.ShowError(err, s.win)
		return
	}
	if len(s.vbox.Objects) == 0 {
		s.win.SetContent(widget.NewLabel("No tasks"))
	} else {
		border := container.NewBorder(s.vbox, nil, nil, nil, nil)
		s.win.SetContent(border)
	}
}
func (s *StatusWindow) Show() {
	s.win.Show()
}
func (s *StatusWindow) ShowAndRun() {
	s.win.ShowAndRun()
}

func main() {
	app := app.New()
	statusWindow := NewStatusWindow(app)
	if desk, ok := app.(desktop.App); ok {
		m := fyne.NewMenu("MyApp",
			fyne.NewMenuItem("Show", func() {
				statusWindow.Show()
			}),
			fyne.NewMenuItem("Change Icon", func() {
				r, err := fyne.LoadResourceFromPath(IconPath(state.StateLowRisk))
				if err != nil {
					panic(err)
				}
				desk.SetSystemTrayIcon(r)
			}))
		r, err := fyne.LoadResourceFromPath("../../resources/Upload.svg")
		if err != nil {
			panic(err)
		}
		desk.SetSystemTrayIcon(r)
		desk.SetSystemTrayMenu(m)
	}
	statusWindow.Update()
	statusWindow.ShowAndRun()
}

func IconPath(s state.State) string {
	return fmt.Sprintf("../../resources/%s.svg", s.String())
}
