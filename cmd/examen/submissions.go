package main

import (
	"examen/pkg/logging"
	"examen/pkg/state"
	"examen/pkg/task"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

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

type SubmissionsWindow struct {
	stopUpdate chan struct{}
	win        fyne.Window
	vbox       *fyne.Container
	from       int
	count      int
	icons      []fyne.Resource
}

func NewSubmissionsWindow(app fyne.App) *SubmissionsWindow {
	s := &SubmissionsWindow{
		stopUpdate: make(chan struct{}),
		win:        app.NewWindow("Submissions"),
		from:       0,
		count:      10,
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
		//logging.Debugf("XXX Close")
		s.Hide()
	})
	return s
}

func (s *SubmissionsWindow) Update() {
	s.vbox.RemoveAll()
	//b := newContextMenuLable("context", m)
	//s.vbox.Add(b)
	//s.list.Add(task.NewTask("C:\\asd\\asd.txt"))
	task.List().IterateIDs(s.from, s.count, func(id task.ID) {
		tsk := task.List().Get(id)
		logging.Debugf("Got id: %v", id)
		downloadItem := fyne.NewMenuItem("Download Report", func() {})
		if tsk.State != state.StateHighRisk &&
			tsk.State != state.StateMediumRisk &&
			tsk.State != state.StateLowRisk {
			downloadItem.Disabled = true
		}
		m := fyne.NewMenu("Examen",
			downloadItem,
			fyne.NewMenuItem("Download Investigation Package", func() {}),
			fyne.NewMenuItem("Delete", func() {
				task.Delete(id)
			}))
		logging.Debugf("Got id: %v, menu: %v", id, m)
		pathLabel := newContextMenuLable(tsk.Path, m)
		logging.Debugf("Got id: %v, pathLabel: %v", id, pathLabel)
		//var icon *tappableIcon
		icon := newTappableIcon(s.icons[task.List().Get(id).State], func() {
			// pup up menu
		})
		stateLabel := widget.NewLabel(tsk.Message)
		logging.Debugf("Got id: %v, icon: %v", id, icon)
		line := container.NewBorder(nil, nil, container.NewHBox(icon, stateLabel, pathLabel), nil)
		logging.Debugf("Add: %v", line)
		s.vbox.Add(line)
	})
	if len(s.vbox.Objects) == 0 {
		s.win.SetContent(widget.NewLabel("No tasks"))
	} else {
		border := container.NewBorder(s.vbox, nil, nil, nil, nil)
		s.win.SetContent(border)
	}
}

func (s *SubmissionsWindow) Show() {
	s.win.Show()
	go func() {
		for {
			logging.Debugf("XXX SubmissionsWindow) Show()")
			select {
			case <-s.stopUpdate:
				logging.Debugf("XXX SubmissionsWindow) Show() Stop")
				return
			case <-task.List().Changes():
				logging.Debugf("XXX SubmissionsWindow) Show() Update")
				s.Update()
			}
		}
	}()
}

func (s *SubmissionsWindow) Hide() {
	s.stopUpdate <- struct{}{}
	s.win.Hide()
}
