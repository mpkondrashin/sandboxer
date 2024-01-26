package main

import (
	"examen/pkg/logging"
	"examen/pkg/state"
	"examen/pkg/task"
	"fmt"
	"image/color"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type SubmissionsWindow struct {
	stopUpdate chan struct{}
	win        fyne.Window
	//vbox       *fyne.Container
	from  int
	count int
	icons []fyne.Resource
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
	//s.vbox = container.NewVBox()
	s.win.SetCloseIntercept(func() {
		//logging.Debugf("XXX Close")
		s.Hide()
	})
	s.win.Resize(fyne.Size{Width: 400, Height: 300})
	//s.win.Content().MinSize (fyne.Size{Width: 400, Height: 300})
	return s
}

func StateColor(s state.State) color.Color {
	switch s {
	//case state.StateUnknown:
	//case state.StateNew:
	//case state.StateUpload:
	//case state.StateInspect:
	//case state.StateReport:
	//return color.RGBA{0, 0, 0, 255}
	case state.StateUnsupported:
		return color.RGBA{158, 158, 158, 255}
	case state.StateError:
		return color.RGBA{255, 0, 0, 255}
	case state.StateNoRisk:
		return color.RGBA{0, 255, 0, 255}
	case state.StateLowRisk:
		return color.RGBA{255, 153, 0, 255}
	case state.StateMediumRisk:
		return color.RGBA{230, 102, 0, 255}
	case state.StateHighRisk:
		return color.RGBA{204, 51, 0, 255}
	default:
		return color.RGBA{0, 0, 0, 255}
	}
}

func (s *SubmissionsWindow) PopUpMenu(tsk *task.Task) *fyne.Menu {
	downloadItem := fyne.NewMenuItem("Download Report", func() {})
	downloadInvestigation := fyne.NewMenuItem("Download Investigation Package", func() {})
	if tsk.State != state.StateHighRisk &&
		tsk.State != state.StateMediumRisk &&
		tsk.State != state.StateLowRisk {
		downloadItem.Disabled = true
		downloadInvestigation.Disabled = true
	}
	var deleteFileItem *fyne.MenuItem
	deleteFileAction := func() {
		dialog.ShowConfirm("Delete file",
			fmt.Sprintf("Following file will be deleted: %s", tsk.Path), func(yes bool) {
				err := os.Remove(tsk.Path)
				if err != nil {
					dialog.ShowError(err, s.win)
					return
				}
				deleteFileItem.Disabled = true
			},
			s.win)
	}

	deleteFileItem = fyne.NewMenuItem("Delete File", deleteFileAction)

	return fyne.NewMenu("Examen",
		downloadItem,
		downloadInvestigation,
		fyne.NewMenuItem("Delete Task", func() {
			task.Delete(tsk.Number)
		}),
		deleteFileItem)
}

func (s *SubmissionsWindow) CardWidget(tsk *task.Task) fyne.CanvasObject {
	path := tsk.Path
	uri := storage.NewFileURI(path)
	icon := container.NewPadded(widget.NewFileIcon(uri))
	//icon.Resize(fyne.Size{Width: 100, Height: 100})
	clr := StateColor(tsk.State)
	fileNameText := canvas.NewText(filepath.Base(path), color.Black)
	fileNameText.TextStyle = fyne.TextStyle{Bold: true}
	stateText := canvas.NewText(tsk.State.String(), clr)
	//stateText.Color = clr
	messageText := canvas.NewText(tsk.Message, clr)
	messageText.TextStyle = fyne.TextStyle{Italic: true}
	messageText.TextSize = 10
	//messageText.Color = StateColor(tsk.State)
	stateVBox := container.NewHBox(stateText, messageText)
	vbox := container.NewPadded(container.NewVBox(
		fileNameText,
		stateVBox,
	))
	t := theme.DefaultTheme()
	menuIcon := newContextMenuIcon(
		t.Icon(theme.IconNameMoreVertical),
		s.PopUpMenu(tsk),
	)
	return container.NewBorder(
		nil, nil, icon, menuIcon,
		vbox,
	)
}

func (s *SubmissionsWindow) Update() {
	//s.vbox.RemoveAll()
	//b := newContextMenuLable("context", m)
	//s.vbox.Add(b)
	//s.list.Add(task.NewTask("C:\\asd\\asd.txt"))
	var list *widget.List
	task.List().Process(func(ids []task.ID) {
		if len(ids) == 0 {
			return
		}
		list = widget.NewList(
			func() int {
				return len(ids)
			},
			func() fyne.CanvasObject {
				//icon := newTappableIcon(s.icons[state.StateHighRisk], nil)
				return s.CardWidget(task.NewTask("placeholder")) /*container.NewPadded(container.NewBorder(
				widget.NewLabel("message"), nil, icon, nil,
				widget.NewLabel("message")))*/
			},
			func(i widget.ListItemID, o fyne.CanvasObject) {
				padded := o.(*fyne.Container)
				padded.RemoveAll()
				tsk := task.List().Get(ids[i])
				//container := container.NewStack(border)
				card := s.CardWidget(tsk)
				padded.Add(card)
			})
		/*
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
			icon := newTappableIcon(s.icons[task.List().Get(id).State], func() {
				// pup up menu
			})
			stateLabel := widget.NewLabel(tsk.Message)
			logging.Debugf("Got id: %v, icon: %v", id, icon)
			line := container.NewBorder(nil, nil, icon, nil, container.NewVBox(pathLabel, stateLabel))
			//logging.Debugf("Add: %v", line)
			s.vbox.Add(line)*/
	})
	if list != nil {
		s.win.SetContent(list)
	} else {
		s.win.SetContent(widget.NewLabel("No submissions"))
	}
	/*if len(s.vbox.Objects) == 0 {
		s.win.SetContent(widget.NewLabel("No tasks"))
	} else {
		border := container.NewBorder(s.vbox, nil, nil, nil, nil)
		s.win.SetContent(border)
	}*/
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
