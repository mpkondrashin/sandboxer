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
	}
	s.cardsList = widget.NewList(
		s.CardsListLength,
		s.CardsListCreateItem,
		s.CardsListUpdateItem,
	)
	s.buttonPrev = widget.NewButton(" < ", s.Prev)
	s.buttonPrev.Disable()
	s.buttonNext = widget.NewButton(" > ", s.Next)
	s.buttonNext.Disable()

	s.pageLabel.TextSize = 12
	s.PopulateOnScreenTasks()
	return s
}

func (s *SubmissionsWindow) Name() string {
	return "Submissions"
}

func (s *SubmissionsWindow) Icon() fyne.Resource {
	return theme.StorageIcon()
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
	//tsk := task.NewTask(0, task.FileTask, "placeholder")
	return s.CardWidget()
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
	reportItem := fyne.NewMenuItem("Show Report", func() {
		s.RunOpen(tsk.Report)
	})
	reportItem.Disabled = tsk.Report == ""
	reportItem.Icon = theme.BrokenImageIcon()

	investigationItem := fyne.NewMenuItem("Investigation Package", func() {
		s.OpenInvestigation(tsk.Investigation)
	})
	investigationItem.Disabled = tsk.Investigation == ""
	investigationItem.Icon = theme.ListIcon()

	var deleteFileItem *fyne.MenuItem
	deleteFileAction := func() {
		dialog.ShowConfirm("Delete file",
			fmt.Sprintf("Following file will be deleted:\n%s", tsk.Path), func(yes bool) {
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
	deleteFileItem.Disabled = !tsk.RiskLevel.IsThreat()
	deleteFileItem.Icon = theme.DeleteIcon()

	recheckAction := func() {
		tsk.SetMessage("")
		tsk.SetRiskLevel(sandbox.RiskLevelUnknown)
		tsk.SetChannel(task.ChPrefilter)
		s.channels.TaskChannel[task.ChPrefilter] <- tsk.Number
	}
	recheckItem := fyne.NewMenuItem("Recheck File", recheckAction)
	recheckItem.Icon = theme.SearchReplaceIcon()

	deleteTaskItem := fyne.NewMenuItem("This Task", func() {
		s.DeleteTask(tsk)
	})
	deleteSameTasksItem := fyne.NewMenuItem("All "+tsk.RiskLevel.String()+" Tasks", func() {
		s.DeleteSameTasks(tsk)
	})
	deleteSameTasksItem.Disabled = tsk.Channel != task.ChDone
	deleteAllTasksItem := fyne.NewMenuItem("All Tasks", s.DeleteAllTasks)

	deleteItem := fyne.NewMenuItem("Delete", nil)
	deleteItem.Icon = theme.CancelIcon()
	deleteItem.ChildMenu = fyne.NewMenu(globals.AppName,
		deleteTaskItem,
		deleteSameTasksItem,
		deleteAllTasksItem,
	)

	return fyne.NewMenu(globals.AppName,
		reportItem,
		investigationItem,
		recheckItem,
		deleteItem,
		deleteFileItem)
}

func (s *SubmissionsWindow) DeleteTask(tsk *task.Task) {
	err := s.list.DeleteTask(tsk)
	if err != nil {
		dialog.ShowError(err, s.win)
		logging.LogError(err)
		return
	}
}

func (s *SubmissionsWindow) DeleteSameTasks(tsk *task.Task) {
	err := s.list.DeleteSameTasks(tsk)
	if err != nil {
		dialog.ShowError(err, s.win)
		logging.LogError(err)
		return
	}
}

func (s *SubmissionsWindow) DeleteAllTasks() {
	err := s.list.DeleteAllTasks()
	if err != nil {
		dialog.ShowError(err, s.win)
		logging.LogError(err)
		return
	}
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

func IconWidget() fyne.CanvasObject {
	iconData := fileicon.VanilaIcon()
	iconResource := &fyne.StaticResource{
		StaticName:    "some file path",
		StaticContent: []byte(iconData),
	}
	image := canvas.NewImageFromResource(iconResource)
	image.SetMinSize(fyne.NewSize(26, 26))
	image.FillMode = canvas.ImageFillContain
	t := canvas.NewText(".", color.RGBA{128, 128, 128, 255})
	t.TextStyle = fyne.TextStyle{Bold: true}
	t.TextSize = 10.0
	labelBorder := container.NewCenter(t)
	return container.NewStack(image, labelBorder)
}

func (s *SubmissionsWindow) CardWidget() fyne.CanvasObject {
	path := "Example Of File.some extension"
	icon := IconWidget()
	fileNameText := canvas.NewText(path, color.Black)

	stateText := canvas.NewText(task.ChResult.String(), color.Black)
	stateText.TextStyle = fyne.TextStyle{Bold: false}

	messageText := canvas.NewText("Unsupported file type", color.Black)
	messageText.TextStyle = fyne.TextStyle{Italic: true}
	messageText.TextSize = 10

	stateVBox := container.NewHBox(stateText, messageText)
	vbox := container.NewPadded(container.NewVBox(
		fileNameText,
		stateVBox,
	))
	menuIcon := newContextMenuIcon(
		theme.DefaultTheme().Icon(theme.IconNameMoreVertical),
		nil,
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
