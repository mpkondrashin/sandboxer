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
	ModalWindow
	//hidden     bool
	stopUpdate chan struct{}
	//win                       fyne.Window
	//enableSubmissionsMenuItem func()
	vbox  *fyne.Container
	from  int
	count int
	//icons []fyne.Resource
	list     *task.TaskList
	channels *dispatchers.Channels
}

func NewSubmissionsWindow(modalWindow ModalWindow, channels *dispatchers.Channels, list *task.TaskList) *SubmissionsWindow {
	s := &SubmissionsWindow{
		ModalWindow: modalWindow,
		//hidden:     true,
		stopUpdate: make(chan struct{}),
		//win:                       app.NewWindow("Submissions"),
		//enableSubmissionsMenuItem: enableSubmissionsMenuItem,
		from:     0,
		count:    10,
		list:     list,
		channels: channels,
		vbox:     container.NewVBox(widget.NewLabel("No Sumbissions")),
	}
	//	s.win.SetCloseIntercept(func() {
	//		s.Hide()
	//	})
	s.win.Resize(fyne.Size{Width: 400, Height: 300})
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
	//	if !s.hidden {
	//		return
	//}
	//s.hidden = false
	s.win.Show()
	fps := time.Millisecond * 300
	go func() {
		haveChanges := true
		for {
			select {
			case <-s.stopUpdate:
				s.Hide() // s.enableSubmissionsMenuItem()
				return
			case <-s.list.Changes():
				haveChanges = true
			case <-time.After(fps):
				if !haveChanges {
					break
				}
				logging.Debugf("Update")
				s.Update()
				haveChanges = false
			}
		}
	}()
}

func (s *SubmissionsWindow) Hide() {
	//s.hidden = true
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
