package main

import (
	"fmt"
	"image/color"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/mpkondrashin/fileicon"

	"sandboxer/pkg/dispatchers"
	"sandboxer/pkg/globals"
	"sandboxer/pkg/logging"
	"sandboxer/pkg/task"
)

type SubmissionsWindow struct {
	ModalWindow
	stopUpdate chan struct{}
	//enableSubmissionsMenuItem func()
	vbox *fyne.Container
	//pageLabel *widget.Label
	pageLabel *canvas.Text
	from      int
	count     int
	//icons []fyne.Resource
	list     *task.TaskList
	channels *dispatchers.Channels
}

func NewSubmissionsWindow(modalWindow ModalWindow, channels *dispatchers.Channels, list *task.TaskList) *SubmissionsWindow {
	s := &SubmissionsWindow{
		ModalWindow: modalWindow,
		stopUpdate:  make(chan struct{}),
		//win:                       app.NewWindow("Submissions"),
		//enableSubmissionsMenuItem: enableSubmissionsMenuItem,
		from:      0,
		count:     10,
		pageLabel: canvas.NewText("", color.Black),
		list:      list,
		channels:  channels,
		vbox:      container.NewVBox(widget.NewLabel("No Sumbissions")),
		//vbox:     container.NewVBox(widget.NewLabel("No Sumbissions")),
	}
	s.pageLabel.TextSize = 12
	//stateText := canvas.NewText(tsk.GetState(), tsk.RiskLevel.Color())

	s.ModalWindow.SetQuit(func() {
		logging.Debugf("stopUpdate")
		s.stopUpdate <- struct{}{}
	})
	//f := widget.NewLabel("a")

	s.win.SetContent(s.Content())
	//	s.win.SetCloseIntercept(func() {
	//		s.Hide()
	//	})
	s.win.Resize(fyne.Size{Width: 400, Height: 300})
	//logging.Debugf("s = %v", s.ModalWindow)
	return s
}

func (s *SubmissionsWindow) Content() fyne.CanvasObject {
	navigationHBox := container.NewHBox(
		widget.NewButton("<", s.Prev),
		widget.NewButton(">", s.Next),
	)
	buttons := container.NewBorder(
		nil,
		nil,
		nil,
		navigationHBox,
	)
	return container.NewBorder(
		s.pageLabel,
		buttons,
		nil,
		nil,
		container.NewScroll(s.vbox),
	)
}
func (s *SubmissionsWindow) Next() {
	l := s.list.Length()
	if s.from+s.count >= l {
		return
	}
	s.from += s.count
	s.Update()
}
func (s *SubmissionsWindow) Prev() {
	s.from -= s.count
	if s.from < 0 {
		s.from = 0
	}
	//f.SetText(fmt.Sprintf("%d", s.from))
	s.Update()
}

func (s *SubmissionsWindow) PopUpMenu(tsk *task.Task) *fyne.Menu {
	downloadItem := fyne.NewMenuItem("Show Report", func() {
		s.OpenReport(tsk.Report)
	})
	downloadItem.Disabled = tsk.Report == ""
	downloadInvestigation := fyne.NewMenuItem("Download Investigation Package", func() {})
	if tsk.State != task.StateDone {
		//		downloadItem.Disabled = true
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
		tsk.SetState(task.StateNew)
		s.channels.TaskChannel[dispatchers.ChPrefilter] <- tsk.Number
	}
	recheckItem := fyne.NewMenuItem("Recheck File", recheckAction)
	if tsk.RiskLevel != task.RiskLevelError {
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

func (s *SubmissionsWindow) OpenReport(report string) {
	cmd := exec.Command("open", report)
	err := cmd.Run()
	if err != nil {
		dialog.ShowError(err, s.win)
	}
}

func IconForFile(path string) fyne.CanvasObject {
	iconData, err := fileicon.FileIcon(path)
	if err != nil {
		iconData = fileicon.VanilaIcon()
	}
	iconResource := &fyne.StaticResource{
		StaticName:    filepath.Base(path),
		StaticContent: []byte(iconData),
	}
	icon := canvas.NewImageFromResource(iconResource)
	icon.SetMinSize(fyne.NewSize(26, 26))
	icon.FillMode = canvas.ImageFillContain
	if err == nil {
		return icon
	}
	ext := filepath.Ext(path)
	if len(ext) > 0 {
		ext = ext[1:]
	}
	var maxSize float32 = 10
	var size float32
	if len(ext) > 0 {
		size = maxSize * 3 / float32(len(ext))
	}
	if size > maxSize {
		size = maxSize
	}
	t := canvas.NewText(ext, color.RGBA{128, 128, 128, 255})
	t.TextStyle = fyne.TextStyle{
		Bold: true,
	}
	t.TextSize = size
	labelBorder := container.NewCenter(t)
	return container.NewStack(icon, labelBorder)
}

func (s *SubmissionsWindow) CardWidget(tsk *task.Task) fyne.CanvasObject {
	path := tsk.Path
	icon := IconForFile(path)
	//uri := storage.NewFileURI(path)
	//icon := container.NewPadded(widget.NewFileIcon(uri))
	//icon.Resize(fyne.Size{Width: 100, Height: 100})
	fileNameText := canvas.NewText(filepath.Base(path), color.Black)
	fileNameText.TextStyle = fyne.TextStyle{Bold: true}
	stateText := canvas.NewText(tsk.GetState(), tsk.RiskLevel.Color())
	//stateText.Color = clr
	//logging.Debugf("XXX MESSAGE GET: %v", tsk)
	messageText := canvas.NewText(tsk.Message, tsk.RiskLevel.Color())
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
		nil, nil, container.NewPadded(icon), menuIcon,
		vbox,
	)
}

func (s *SubmissionsWindow) Update() {
	s.pageLabel.Text = fmt.Sprintf("Submissions %d - %d out of %d", s.from+1, s.from+s.count+1, s.list.Length())
	s.pageLabel.Refresh()
	//	logging.Debugf("XXX SubmissionsWindow.Update()")
	//s.l.SetText(s.l.Text + "!")
	//	return
	s.vbox.RemoveAll()
	//b := newContextMenuLable("context", m)
	//s.vbox.Add(b)
	//s.list.Add(task.NewTask("C:\\asd\\asd.txt"))

	s.list.Process(func(ids []task.ID) {
		//logging.Debugf("XXX SubmissionsWindow.Update() Process")
		//for i, idx := range ids {
		for i := s.from; i < s.from+s.count && i < len(ids); i++ {
			idx := ids[i]
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
		//s.win.SetContent(container.NewScroll(s.vbox))
	} else {
		s.vbox.Add(widget.NewLabel("No submissions"))
		//s.win.SetContent(widget.NewLabel("No submissions"))
	}
	s.vbox.Refresh()
	//s.win.Canvas().Refresh()
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
				logging.Debugf("Stop Update")
				//				s.Hide() // s.enableSubmissionsMenuItem()
				return
			case <-s.list.Changes():
				haveChanges = true
			case <-time.After(fps):
				if !haveChanges {
					break
				}
				logging.Debugf(strings.Repeat("*", 100))
				s.Update()
				haveChanges = false
			}
		}
	}()
}

/*
func (s *SubmissionsWindow) Hide() {
	//s.hidden = true

	s.ModalWindow.Hide() // s.win.Hide()
}
*/
