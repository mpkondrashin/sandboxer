package main

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"sandboxer/pkg/dispatchers"
	"sandboxer/pkg/globals"
	"sandboxer/pkg/logging"
	"sandboxer/pkg/state"
	"sandboxer/pkg/task"
)

type SubmissionsWindow struct {
	stopUpdate chan struct{}
	win        fyne.Window
	vbox       *fyne.Container
	from       int
	count      int
	//icons []fyne.Resource
	list     *task.TaskList
	channels *dispatchers.Channels
}

func NewSubmissionsWindow(app fyne.App, channels *dispatchers.Channels, list *task.TaskList) *SubmissionsWindow {
	s := &SubmissionsWindow{
		stopUpdate: make(chan struct{}),
		win:        app.NewWindow("Submissions"),
		from:       0,
		count:      10,
		list:       list,
		channels:   channels,
		vbox:       container.NewVBox(widget.NewLabel("No Sumbissions")),
	}
	//	for st := state.StateNew; st < state.StateCount; st++ {
	//		r, err := fyne.LoadResourceFromPath(IconPath(st))
	//		if err != nil {
	//			panic(err)
	//		}
	//		s.icons = append(s.icons, r)
	//	}
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
	case state.StateUnsupported, state.StateIgnored:
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
				if !yes {
					return
				}
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
	recheckAction := func() {
		tsk.SetState(state.StateNew)
		s.channels.TaskChannel[dispatchers.ChPrefilter] <- tsk.Number
	}
	recheckItem := fyne.NewMenuItem("Recheck File", recheckAction)
	if tsk.State != state.StateError {
		recheckItem.Disabled = true
	}
	return fyne.NewMenu(globals.AppName,
		downloadItem,
		downloadInvestigation,
		recheckItem,
		fyne.NewMenuItem("Delete Task", func() {
			s.list.DelByID(tsk.Number)
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
	stateText := canvas.NewText(Split(tsk.State.String()), clr)
	//stateText.Color = clr
	logging.Debugf("XXX MESSAGE GET: %v", tsk)
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

func (s *SubmissionsWindow) _Update() {
	logging.Debugf("XXX SubmissionsWindow.Update()")
	//s.vbox.RemoveAll()
	//b := newContextMenuLable("context", m)
	//s.vbox.Add(b)
	//s.list.Add(task.NewTask("C:\\asd\\asd.txt"))
	var list *widget.List
	s.list.Process(func(ids []task.ID) {
		logging.Debugf("XXX SubmissionsWindow.Update() Process")
		if len(ids) == 0 {
			return
		}
		idsCopy := ids[:]
		list = widget.NewList(
			func() int {
				return len(idsCopy)
			},
			func() fyne.CanvasObject {
				//icon := newTappableIcon(s.icons[state.StateHighRisk], nil)
				return s.CardWidget(task.NewTask(0, "placeholder")) /*container.NewPadded(container.NewBorder(
				widget.NewLabel("message"), nil, icon, nil,
				widget.NewLabel("message")))*/
			},
			func(i widget.ListItemID, o fyne.CanvasObject) {
				padded := o.(*fyne.Container)
				padded.RemoveAll()
				_ = s.list.Task(idsCopy[i], func(tsk *task.Task) error {
					if tsk == nil {
						tsk = task.NewTask(0, "placeholder")
						//logging.Debugf("tsk = nil, i = %d, ids[i]=%d , ids = %v", i, ids[i], ids)
					}
					card := s.CardWidget(tsk)
					padded.Add(card)
					return nil
				})
			})
	})
	if list != nil {
		logging.Debugf("XXX SubmissionsWindow.Update() List.Length = %d", list.Length())
		s.win.SetContent(list)
	} else {
		s.win.SetContent(widget.NewLabel("No submissions"))
	}
}

func (s *SubmissionsWindow) Update() {
	logging.Debugf("XXX SubmissionsWindow.Update()")
	s.vbox.RemoveAll()
	//b := newContextMenuLable("context", m)
	//s.vbox.Add(b)
	//s.list.Add(task.NewTask("C:\\asd\\asd.txt"))
	s.list.Process(func(ids []task.ID) {
		//logging.Debugf("XXX SubmissionsWindow.Update() Process")
		for i, idx := range ids {
			_ = s.list.Task(idx, func(tsk *task.Task) error {
				if tsk == nil {
					tsk = task.NewTask(0, "placeholder")
					//logging.Debugf("tsk = nil, i = %d, ids[i]=%d , ids = %v", i, ids[i], ids)
				}
				card := s.CardWidget(tsk)
				//				if i > 0 {
				//				}
				_ = i
				s.vbox.Add(card) // padded
				s.vbox.Add(canvas.NewLine(color.RGBA{158, 158, 158, 255}))
				return nil
			})
		}
	})
	if len(s.vbox.Objects) > 0 {
		//logging.Debugf("XXX SubmissionsWindow.Update() List.Length = %d", len(s.vbox.Objects))
		s.win.SetContent(container.NewScroll(s.vbox))
	} else {
		s.win.SetContent(widget.NewLabel("No submissions"))
	}
}

func (s *SubmissionsWindow) Show() {
	s.win.Show()
	fps := time.Millisecond * 100
	go func() {
		// updateTime := time.Now().Add(fps)
		haveChanges := true
		for {
			select {
			case <-s.stopUpdate:
				return
			case <-s.list.Changes():
				haveChanges = true
			case <-time.After(fps):
				logging.Debugf("Update")
				if !haveChanges {
					break
				}
				// if time.Now().After(updateTime) {
				s.Update()
				haveChanges = false
				// updateTime = time.Now().Add(fps)
				// }
			}
		}
	}()
}

func (s *SubmissionsWindow) Hide() {
	s.stopUpdate <- struct{}{}
	s.win.Hide()
}

// Split - State output. Take from https://github.com/fatih/camelcase
func Split(src string) string {
	// don't split invalid utf8
	if !utf8.ValidString(src) {
		return src
	}
	entries := []string{}
	var runes [][]rune
	lastClass := 0
	class := 0
	// split into fields based on class of unicode character
	for _, r := range src {
		switch true {
		case unicode.IsLower(r):
			class = 1
		case unicode.IsUpper(r):
			class = 2
		case unicode.IsDigit(r):
			class = 3
		default:
			class = 4
		}
		if class == lastClass {
			runes[len(runes)-1] = append(runes[len(runes)-1], r)
		} else {
			runes = append(runes, []rune{r})
		}
		lastClass = class
	}
	// handle upper case -> lower case sequences, e.g.
	// "PDFL", "oader" -> "PDF", "Loader"
	for i := 0; i < len(runes)-1; i++ {
		if unicode.IsUpper(runes[i][0]) && unicode.IsLower(runes[i+1][0]) {
			runes[i+1] = append([]rune{runes[i][len(runes[i])-1]}, runes[i+1]...)
			runes[i] = runes[i][:len(runes[i])-1]
		}
	}
	// construct []string from results
	for _, s := range runes {
		if len(s) > 0 {
			entries = append(entries, string(s))
		}
	}
	return strings.Join(entries, " ")
}
