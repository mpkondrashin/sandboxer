package main

import (
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"examen/pkg/shm"
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

func main() {
	/*s := state.NewMemState()
	for st := state.StateUnknown; st < state.StateCount; st++ {
		fileName := fmt.Sprintf("C:\\Documents\\file_%v.docx", st)
		id := fmt.Sprint(st)
		s.AddObject(state.NewObject(id, fileName))
		s.SetState(id, st)
	}
	fmt.Println(s.ListObjects())*/

	app := app.New()
	win := app.NewWindow("Status")

	reader, err := shm.NewSHMReader(10000)
	if err != nil {
		dialog.ShowError(err, win)
		return
	}
	//list, err := s.ListObjects()
	//_ = err
	var icons []fyne.Resource
	for s := state.StateUnknown; s < state.StateCount; s++ {
		r, err := fyne.LoadResourceFromPath(IconPath(s))
		if err != nil {
			panic(err)
		}
		icons = append(icons, r)
	}
	vbox := container.NewVBox()
	l := task.NewList()
	err = reader.Read(l)
	if err != nil {
		dialog.ShowError(err, win)
		return
	}
	l.Iterate(func(t *task.Task) {
		//iconLabel := widget.NewLabel(o.State.String())
		pathLabel := widget.NewLabel(t.Path)
		var icon *tappableIcon
		icon = newTappableIcon(icons[t.State], func() {
			go func() {
				for s := state.StateUnknown; s <= state.StateHighRisk; s++ {
					icon.SetResource(icons[s])
					time.Sleep(1 * time.Second)
				}
			}()
		})
		line := container.NewBorder(nil, nil, container.NewHBox(icon, pathLabel), nil)
		vbox.Add(line)
	})
	if len(vbox.Objects) == 0 {
		win.SetContent(widget.NewLabel("No tasks"))
	} else {
		border := container.NewBorder(vbox, nil, nil, nil, nil)
		win.SetContent(border)
	}
	//win.Resize(fyne.NewSize(600, 400))
	win.ShowAndRun()
}

func IconPath(s state.State) string {
	return fmt.Sprintf("../../resources/%s.svg", s.String())
}
