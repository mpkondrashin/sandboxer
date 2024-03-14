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
	"path/filepath"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/lutzky/go-bidi"

	"github.com/mpkondrashin/fileicon"

	"sandboxer/pkg/config"
	"sandboxer/pkg/globals"
	"sandboxer/pkg/logging"
	"sandboxer/pkg/sandbox"
	"sandboxer/pkg/task"
	"sandboxer/pkg/xplatform"
)

type SubmissionsWindow struct {
	mx            sync.Mutex
	win           fyne.Window
	stopUpdate    chan struct{}
	conf          *config.Configuration
	cardsList     *widget.List
	buttonNext    *widget.Button
	buttonPrev    *widget.Button
	pageLabel     *canvas.Text
	statusLabel   *widget.Label
	from          int
	count         int
	onScreenTasks []*task.Task
	list          *task.TaskList
	channels      *task.Channels
}

func NewSubmissionsWindow(channels *task.Channels, list *task.TaskList, conf *config.Configuration) *SubmissionsWindow {
	s := &SubmissionsWindow{
		stopUpdate:  make(chan struct{}),
		conf:        conf,
		from:        0,
		count:       10,
		pageLabel:   canvas.NewText("", color.Black),
		statusLabel: widget.NewLabel(""),
		list:        list,
		channels:    channels,
		//vbox:        container.NewVBox(widget.NewLabel("No Sumbissions")),
	}
	s.cardsList = widget.NewList(
		s.CardsListLength,
		s.CardsListCreateItem,
		s.CardsListUpdateItem,
	)
	//s.cardsList.OnSelected = func(id widget.ListItemID) {
	//	menu := s.PopUpMenu(task.NewTask(0, task.FileTask, "placeholder"))
	//s.cardsList.
	//	widget.ShowPopUpMenuAtPosition(menu, fyne.CurrentApp().Driver().CanvasForObject(s.cardsList), fyne.NewPos(0, 0)) // e.AbsolutePosition)

	//	fmt.Println("Onselected", id)
	//}
	s.buttonPrev = widget.NewButton("<", s.Prev)
	s.buttonPrev.Disable()
	s.buttonNext = widget.NewButton(">", s.Next)
	s.buttonNext.Disable()

	//TMPs.pageLabel.TextSize = 12
	s.PopulateOnScreenTasks()
	return s
}

func (s *SubmissionsWindow) Name() string {
	return "Submissions"
}

func (s *SubmissionsWindow) Content(w *ModalWindow) fyne.CanvasObject {
	s.win = w.win
	s.win.Resize(fyne.Size{Width: 400, Height: 300})
	navigationHBox := container.NewHBox(
		s.buttonPrev,
		s.buttonNext,
	)
	buttons := container.NewBorder(
		nil,
		nil,
		s.statusLabel,
		navigationHBox,
	)
	return container.NewBorder(
		s.pageLabel,
		buttons,
		nil,
		nil,
		s.cardsList, //container.NewScroll(s.vbox),
	)
}

func (s *SubmissionsWindow) PopulateOnScreenTasks() {
	s.onScreenTasks = nil
	s.list.Process(func(ids []task.ID) {
		for i := s.from; i < min(s.from+s.count, len(ids)); i++ {
			idx := ids[i]
			_ = s.list.Task(idx, func(tsk *task.Task) error {
				s.onScreenTasks = append(s.onScreenTasks, tsk)
				return nil
			})
		}
	})

	to := s.from + s.count + 1
	if to > s.list.Length() {
		to = s.list.Length()
	}

	if s.from > 0 || s.from+s.count < s.list.Length() {
		s.pageLabel.Text = fmt.Sprintf("Submissions %d - %d out of %d", s.from+1, to, s.list.Length())
		//s.pageLabel.SetText(fmt.Sprintf("Submissions %d - %d out of %d", s.from+1, to, s.list.Length()))
	} else {
		//s.pageLabel.SetText("")
		s.pageLabel.Text = ""
	}
	s.pageLabel.Refresh()

	count := s.list.CountActiveTasks()
	activeTasks := "No active task"
	if count > 0 {
		activeTasks = fmt.Sprintf("Active tasks: %d", s.list.CountActiveTasks())
	}
	s.statusLabel.SetText(activeTasks)
	if s.from > 0 {
		s.buttonPrev.Enable()
	} else {
		s.buttonPrev.Disable()
	}
	if s.from+s.count < s.list.Length() {
		s.buttonNext.Enable()
	} else {
		s.buttonNext.Disable()
	}
	s.cardsList.Refresh()
}
func (s *SubmissionsWindow) Next() {
	l := s.list.Length()
	if s.from+s.count >= l {
		return
	}
	s.from += s.count
	s.PopulateOnScreenTasks()
	// s.Update()
}
func (s *SubmissionsWindow) Prev() {
	s.from -= s.count
	if s.from < 0 {
		s.from = 0
	}
	s.PopulateOnScreenTasks()
	//s.Update()
}

func (s *SubmissionsWindow) CardsListLength() int {
	return len(s.onScreenTasks)
}

func (s *SubmissionsWindow) CardsListCreateItem() fyne.CanvasObject {
	tsk := task.NewTask(0, task.FileTask, "placeholder")
	return s.CardWidget(tsk)
}

func (s *SubmissionsWindow) CardsListUpdateItem(itemID widget.ListItemID, object fyne.CanvasObject) {
	//tsk := s.list.Get(task.ID(itemID))
	s.UpdateWidget(s.onScreenTasks[itemID], object)
}

func (s *SubmissionsWindow) UpdateWidget(tsk *task.Task, object fyne.CanvasObject) {
	hBox := object.(*fyne.Container)
	menuIcon := hBox.Objects[0].(*contextMenuIcon)
	menuIcon.Menu = s.PopUpMenu(tsk)
	iconStack := hBox.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container)
	vbox := hBox.Objects[2].(*fyne.Container).Objects[0].(*fyne.Container)
	fileNameText := vbox.Objects[0].(*canvas.Text)
	stateHBox := vbox.Objects[1].(*fyne.Container)
	stateText := stateHBox.Objects[0].(*canvas.Text)
	messageText := stateHBox.Objects[1].(*canvas.Text)
	messageText.Text = tsk.Message
	messageText.Color = tsk.RiskLevel.Color()

	icon := iconStack.Objects[0].(*canvas.Image)
	label := iconStack.Objects[1].(*fyne.Container).Objects[0].(*canvas.Text)

	iconData, err := fileicon.FileIcon(tsk.Path)
	if err != nil {
		iconData = fileicon.VanilaIcon()
	}
	iconResource := &fyne.StaticResource{
		StaticName:    filepath.Base(tsk.Path),
		StaticContent: []byte(iconData),
	}
	icon.Resource = iconResource
	if err != nil {
		label.Text, label.TextSize = ExtAndSize(tsk.Path)
	} else {
		label.Text = ""
	}
	bidiStr, err := bidi.Display(tsk.Title())
	if err != nil {
		logging.LogError(err)
		bidiStr = tsk.Title()
	}

	fileNameText.Text = bidiStr
	stateText.Text = tsk.GetChannel()
	stateText.Color = tsk.RiskLevel.Color()
	stateText.TextStyle = fyne.TextStyle{Bold: tsk.Active}

	messageText.TextStyle = fyne.TextStyle{Italic: true}
	messageText.TextSize = 10
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
	deleteFileItem.Disabled = tsk.RiskLevel != sandbox.RiskLevelHigh &&
		tsk.RiskLevel != sandbox.RiskLevelMedium &&
		tsk.RiskLevel != sandbox.RiskLevelLow
	recheckAction := func() {
		tsk.SetMessage("")
		tsk.SetRiskLevel(sandbox.RiskLevelUnknown)
		tsk.SetChannel(task.ChPrefilter)
		s.channels.TaskChannel[task.ChPrefilter] <- tsk.Number
	}
	recheckItem := fyne.NewMenuItem("Recheck File", recheckAction)
	//recheckItem.Disabled = (tsk.RiskLevel != sandbox.RiskLevelError) && (tsk.RiskLevel != sandbox.RiskLevelUnsupported)
	return fyne.NewMenu(globals.AppName,
		downloadItem,
		downloadInvestigation,
		recheckItem,
		fyne.NewMenuItem("Delete Task", func() {
			s.DeleteTask(tsk)
		}),
		deleteFileItem)
}

func (s *SubmissionsWindow) DeleteTask(tsk *task.Task) {
	err := tsk.Delete()
	if err != nil {
		dialog.ShowError(err, s.win)
		logging.LogError(err)
		return
	}
	s.list.DelByID(tsk.Number)
	//s.PopulateOnScreenTasks()
}

func (s *SubmissionsWindow) RunOpen(path string) {
	err := xplatform.RunOpen(path)
	if err != nil {
		dialog.ShowError(err, s.win)
	}
}

func (s *SubmissionsWindow) OpenInvestigation(investigation string) {
	if !s.conf.GetShowPasswordHint() {
		s.RunOpen(investigation)
		return
	}
	dialog.ShowConfirm("Hint",
		"Password for archive is \"virus\". Show this note next time?",
		func(yes bool) {
			if !yes {
				s.conf.SetShowPasswordHint(false)
				err := s.conf.Save()
				if err != nil {
					dialog.ShowError(err, s.win)
				}
			}
			s.RunOpen(investigation)
		}, s.win)
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

	bidiStr, err := bidi.Display(tsk.Title())
	if err != nil {
		logging.LogError(err)
		bidiStr = tsk.Title()
	}
	fileNameText := canvas.NewText(bidiStr, color.Black)
	//fileNameText.TextStyle = fyne.TextStyle{Bold: true}
	stateText := canvas.NewText(tsk.GetChannel(), tsk.RiskLevel.Color())
	stateText.TextStyle = fyne.TextStyle{Bold: tsk.Active}

	messageText := canvas.NewText(tsk.Message, tsk.RiskLevel.Color())
	messageText.TextStyle = fyne.TextStyle{Italic: true}
	messageText.TextSize = 10

	stateVBox := container.NewHBox(stateText, messageText)
	vbox := container.NewPadded(container.NewVBox(
		fileNameText,
		stateVBox,
	))
	menuIcon := newContextMenuIcon(
		theme.DefaultTheme().Icon(theme.IconNameMoreVertical),
		s.PopUpMenu(tsk),
	)
	return container.NewHBox(menuIcon, container.NewPadded(icon), vbox)
}

func ExtAndSize(path string) (string, float32) {
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
	return ext, size
}

/*
	func (s *SubmissionsWindow) Update() {
		s.mx.Lock()
		defer s.mx.Unlock()
		to := s.from + s.count + 1
		if to > s.list.Length() {
			to = s.list.Length()
		}
		s.vbox.RemoveAll()
		s.list.Process(func(ids []task.ID) {
			count := s.list.CountActiveTasks()
			activeTasks := "No active task"
			if count > 0 {
				activeTasks = fmt.Sprintf("Active tasks: %d", s.list.CountActiveTasks())
			}
			s.statusLabel.SetText(activeTasks)
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
		if s.from > 0 || s.from+s.count < s.list.Length() {
			s.pageLabel.Text = fmt.Sprintf("Submissions %d - %d out of %d", s.from+1, to, s.list.Length())
		} else {
			s.pageLabel.Text = ""
		}
		if len(s.vbox.Objects) == 0 {
			s.vbox.Add(container.NewCenter(widget.NewLabel("No submissions")))
			s.pageLabel.Text = ""
			s.statusLabel.SetText("")
			s.buttonNext.Disable()
			s.buttonPrev.Disable()
		}
		s.pageLabel.Refresh()
		//s.vbox.Refresh()
	}
*/
func (s *SubmissionsWindow) Show() {
	fps := time.Millisecond * 300
	go func() {
		logging.Debugf("Submissions Show")
		haveChanges := true
		for {
			select {
			case <-s.stopUpdate:
				logging.Debugf("Got Stop Update")
				return
			case <-s.list.Changes():
				haveChanges = true
			case <-time.After(fps):
				if !haveChanges {
					break
				}
				s.PopulateOnScreenTasks()

				haveChanges = false
			}
		}
	}()
}

func (s *SubmissionsWindow) Hide() {
	logging.Debugf("Send Stop Update")
	s.stopUpdate <- struct{}{}
}
