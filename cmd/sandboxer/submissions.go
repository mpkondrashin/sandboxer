/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

submissions.go

Submissions list window
*/
package main

import (
	"fmt"
	"image/color"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/mpkondrashin/fileicon"

	"sandboxer/pkg/config"
	"sandboxer/pkg/dispatchers"
	"sandboxer/pkg/globals"
	"sandboxer/pkg/logging"
	"sandboxer/pkg/task"
)

type SubmissionsWindow struct {
	ModalWindow
	stopUpdate chan struct{}
	//enableSubmissionsMenuItem func()
	conf       *config.Configuration
	vbox       *fyne.Container
	buttonNext *widget.Button
	buttonPrev *widget.Button
	//pageLabel *widget.Label
	pageLabel *canvas.Text
	from      int
	count     int
	//icons []fyne.Resource
	list     *task.TaskList
	channels *dispatchers.Channels
}

func NewSubmissionsWindow(modalWindow ModalWindow, channels *dispatchers.Channels, list *task.TaskList, conf *config.Configuration) *SubmissionsWindow {
	s := &SubmissionsWindow{
		ModalWindow: modalWindow,
		stopUpdate:  make(chan struct{}),
		conf:        conf,
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
	s.buttonPrev = widget.NewButton("<", s.Prev)
	s.buttonPrev.Disable()
	s.buttonNext = widget.NewButton(">", s.Next)
	s.buttonNext.Disable()

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
		s.buttonPrev,
		s.buttonNext,
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
		s.RunOpen(tsk.Report)
	})
	downloadItem.Disabled = tsk.Report == ""

	downloadInvestigation := fyne.NewMenuItem("Investigation Package", func() {
		s.OpenInvestigation(tsk.Investigation)
	})
	downloadInvestigation.Disabled = tsk.Investigation == ""

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
	deleteFileItem.Disabled = tsk.RiskLevel != task.RiskLevelHigh &&
		tsk.RiskLevel != task.RiskLevelMedium &&
		tsk.RiskLevel != task.RiskLevelLow
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

func (s *SubmissionsWindow) RunOpen(path string) {
	err := RunOpen(path)
	if err != nil {
		dialog.ShowError(err, s.win)
	}
}

func RunOpen(path string) error {
	name := "open"
	args := []string{path}
	if runtime.GOOS == "windows" {
		name = "cmd"
		args = []string{"/C", "start", path}
	}
	cmd := exec.Command(name, args...)
	return cmd.Run()
}

func (s *SubmissionsWindow) OpenInvestigation(investigation string) {
	//logging.Debugf("OpenInvestigation.Showhint: %v", s.conf.ShowPasswordHint)
	if !s.conf.ShowPasswordHint {
		s.RunOpen(investigation)
		return
	}
	dialog.ShowConfirm("Hint",
		"Password for archive is \"virus\". Show this note next time?",
		func(yes bool) {
			//logging.Debugf("OpenInvestigation.Yes: %v", yes)
			if !yes {
				s.conf.ShowPasswordHint = false
				err := s.conf.Save()
				if err != nil {
					dialog.ShowError(err, s.win)
				}
				//	logging.Debugf("OpenInvestigation.Save: %v", err)
			}
			s.RunOpen(investigation)
		}, s.win)
	//logging.Debugf("OpenInvestigation.Done")
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
	return container.NewHBox(menuIcon, container.NewPadded(icon), vbox)
	//container.NewBorder(
	//	nil, nil, container.NewHBox(menuIcon, container.NewPadded(icon)), nil, //menuIcon,
	//		vbox,
	//	)
}

func (s *SubmissionsWindow) Update() {
	to := s.from + s.count + 1
	if to > s.list.Length() {
		to = s.list.Length()
	}
	//	logging.Debugf("XXX SubmissionsWindow.Update()")
	s.vbox.RemoveAll()
	s.list.Process(func(ids []task.ID) {
		for i := s.from; i < s.from+s.count && i < len(ids); i++ {
			idx := ids[i]
			_ = s.list.Task(idx, func(tsk *task.Task) error {
				//if tsk == nil {
				//tsk = task.NewTask(0, "placeholder")
				//}
				card := s.CardWidget(tsk)
				s.vbox.Add(card) // padded
				s.vbox.Add(canvas.NewLine(color.RGBA{158, 158, 158, 255}))
				return nil
			})
		}
		if s.from > 0 {
			s.buttonPrev.Enable()
		} else {
			s.buttonPrev.Disable()
		}
		if s.from+s.count < len(ids) {
			s.buttonNext.Enable()
		} else {
			s.buttonNext.Disable()
		}
	})
	s.pageLabel.Text = fmt.Sprintf("Submissions %d - %d out of %d", s.from+1, to, s.list.Length())
	if len(s.vbox.Objects) == 0 {
		s.vbox.Add(container.NewCenter(widget.NewLabel("No submissions")))
		s.pageLabel.Text = ""
		s.buttonNext.Disable()
		s.buttonPrev.Disable()
	}
	s.pageLabel.Refresh()

	s.vbox.Refresh()
}

func (s *SubmissionsWindow) Show() {
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
				//logging.Debugf(strings.Repeat("*", 100))
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
