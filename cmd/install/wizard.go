/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

wizard.go

Installation wizard
*/
package main

import (
	"errors"
	"fmt"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"golang.org/x/mod/semver"

	"sandboxer/pkg/fatal"
	"sandboxer/pkg/globals"
	"sandboxer/pkg/logging"
)

//go:generate fyne bundle --name IconSVGResource --output resource.go ../../resources/icon.png

var ErrAbort = errors.New("abort")

type Page interface {
	Name() string
	Content() fyne.CanvasObject
	Run()
	AquireData(installer *Installer) error
	Next(previousPage PageIndex) PageIndex
	Prev() PageIndex
}

type BasePage struct {
	wiz          *Wizard
	previousPage PageIndex
}

func NewBasePage(wiz *Wizard) BasePage {
	return BasePage{
		wiz: wiz,
	}
}
func (p *BasePage) SavePrevious(previousPage PageIndex) {
	if previousPage == pgExit {
		return
	}
	p.previousPage = previousPage
}

func (p *BasePage) Prev() PageIndex {
	return p.previousPage
}

func (p *BasePage) Run() {}

func (p *BasePage) AquireData(installer *Installer) error {
	return nil
}

type PageIndex int

const (
	pgIntro PageIndex = iota
	pgDelete
	pgDowngrade
	pgReinstall
	pgUpgrade
	pgProxy
	pgVOSettings
	pgDDSettings
	pgAutostart
	pgFolder
	pgInstallation
	pgUninstall
	pgFinish
	pgExit
)

type Wizard struct {
	pages          []Page
	firstPage      PageIndex
	currentPage    PageIndex
	app            fyne.App
	win            fyne.Window
	pagesList      *fyne.Container
	buttonsLine    *fyne.Container
	installer      *Installer
	capturesFolder string
}

func NewWizard(capturesFolder string) *Wizard {
	installer, err := NewInstaller(globals.AppID)
	if err != nil {
		logging.LogError(err)
		fatal.Warning("Error", err.Error())
		os.Exit(globals.ExitNewInstaller)
	}
	w := &Wizard{
		app:            app.NewWithID(globals.AppID),
		pagesList:      container.NewVBox(),
		buttonsLine:    container.NewHBox(),
		capturesFolder: capturesFolder,
		installer:      installer,
	}
	w.app.Lifecycle()
	w.win = w.app.NewWindow(globals.AppName + " Install Program")
	w.win.Resize(fyne.NewSize(600, 400))
	//c.win.SetFixedSize(true)
	w.win.SetMaster()
	w.firstPage = w.Pages()
	w.currentPage = w.firstPage
	prtScr := &desktop.CustomShortcut{KeyName: fyne.KeyI, Modifier: fyne.KeyModifierControl}
	w.win.Canvas().AddShortcut(prtScr, w.captureWindowContents)
	w.win.SetContent(w.Window())
	logging.Debugf("NewWizard: %v", w)
	return w
}

func (w *Wizard) Pages() PageIndex {
	w.pages = make([]Page, pgExit)
	w.pages[pgIntro] = &PageIntro{BasePage: NewBasePage(w)}
	w.pages[pgDowngrade] = &PageDowngrade{BasePage: NewBasePage(w)}
	w.pages[pgReinstall] = &PageReinstall{BasePage: NewBasePage(w)}
	w.pages[pgUpgrade] = &PageUpgrade{BasePage: NewBasePage(w)}
	w.pages[pgProxy] = &PageProxy{BasePage: NewBasePage(w)}
	w.pages[pgVOSettings] = &PageVOToken{BasePage: NewBasePage(w)}
	w.pages[pgDDSettings] = &PageDDAnSettings{BasePage: NewBasePage(w)}
	w.pages[pgAutostart] = &PageAutostart{BasePage: NewBasePage(w)}
	w.pages[pgFolder] = &PageFolder{BasePage: NewBasePage(w)}
	w.pages[pgInstallation] = &PageInstallation{BasePage: NewBasePage(w)}
	w.pages[pgUninstall] = &PageUninstall{BasePage: NewBasePage(w)}
	w.pages[pgFinish] = &PageFinish{BasePage: NewBasePage(w)}

	err := w.installer.config.Load()
	logging.LogError(err)
	if err != nil {
		if os.IsNotExist(err) {
			return pgIntro
		}
		w.pages[pgDelete] = &PageDelete{
			BasePage:     NewBasePage(w),
			ErrorMessage: err.Error(),
		}
		return pgDelete
	}
	cmp := semver.Compare(globals.Version, w.installer.config.GetVersion())
	logging.Infof("Installer version: %s. Config version %s. Compare %d", globals.Version, w.installer.config.GetVersion(), cmp)
	switch cmp {
	case -1:
		return pgDowngrade
	case 0:
		return pgReinstall
	case 1:
		return pgUpgrade
	}
	return pgIntro
}

func (c *Wizard) captureWindowContents(shortcut fyne.Shortcut) {
	if c.capturesFolder == "" {
		log.Println("--capture is not set")
		return
	}
	fileName := fmt.Sprintf("page_%d.png", c.currentPage)
	filePath := filepath.Join(c.capturesFolder, fileName)
	f, err := os.Create(filePath)
	if err != nil {
		dialog.ShowError(err, c.win)
		return
	}
	defer f.Close()
	image := c.win.Canvas().Capture()
	if err := png.Encode(f, image); err != nil {
		dialog.ShowError(err, c.win)
		return
	}
}

/*
// DELETE
	func (c *NSHIControl) SaveScreenShots() {
		time.Sleep(1 * time.Second)
		for c.current = 0; c.current < len(c.pages); c.current++ {
			c.win.SetContent(c.Window(c.pages[c.current]))
			c.win.Show()
			c.CaptureImage()
			time.Sleep(1 * time.Second)
		}
	}
*/

func (c *Wizard) Window() fyne.CanvasObject {
	logging.Debugf("Window")
	p := c.pages[c.currentPage]
	c.UpdatePagesList()
	middle := container.NewPadded(container.NewVBox(layout.NewSpacer(), p.Content(), layout.NewSpacer()))
	upper := container.NewBorder(nil, nil, container.NewHBox(c.pagesList, widget.NewSeparator()), nil, middle)
	buttons := container.NewBorder(nil, nil, nil, c.buttonsLine)
	bottom := container.NewVBox(widget.NewSeparator(), buttons)
	return container.NewBorder(nil, container.NewPadded(bottom), nil, nil, upper)
}

func (c *Wizard) UpdatePagesList() {
	c.pagesList.RemoveAll()
	image := canvas.NewImageFromResource(ApplicationIcon)
	image.SetMinSize(fyne.NewSize(52, 52))
	image.FillMode = canvas.ImageFillContain
	c.pagesList.Add(image)
	previous := pgExit
	i := c.firstPage
	for {
		if i == pgExit {
			break
		}
		pg := c.pages[i]
		next := pg.Next(previous)
		if i == c.currentPage {
			c.pagesList.Add(widget.NewLabelWithStyle("▶ "+pg.Name(), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
			prev, next := c.Buttons(i == c.firstPage, next == pgExit)
			c.buttonsLine.RemoveAll()
			c.buttonsLine.Add(prev)
			c.buttonsLine.Add(next)
		} else {
			c.pagesList.Add(widget.NewLabel("    " + pg.Name()))
		}
		previous = i
		i = next
	}
}

func (c *Wizard) Buttons(first, last bool) (*widget.Button, *widget.Button) {
	prevButton := widget.NewButtonWithIcon("Back", theme.NavigateBackIcon(), c.Prev)
	if first {
		prevButton.Disable()
	}

	nextButton := widget.NewButtonWithIcon("Next", theme.NavigateNextIcon(), c.Next)
	nextButton.IconPlacement = widget.ButtonIconTrailingText

	if last {
		nextButton = widget.NewButtonWithIcon("Quit", theme.CancelIcon(), c.Quit)
	}
	return prevButton, nextButton
}

func (c *Wizard) Quit() {
	logging.Debugf("Quit")
	err := c.pages[c.currentPage].AquireData(c.installer)
	if err != nil {
		logging.Errorf("AquireData: %v", err)
		dialog.ShowError(err, c.win)
	}
	//dialog.ShowConfirm("Sandboxer", "Exit?", )
	c.app.Quit()
}

func (c *Wizard) Next() {
	logging.Debugf("Next from page %d", c.currentPage)
	err := c.pages[c.currentPage].AquireData(c.installer)
	if err != nil {
		if errors.Is(err, ErrAbort) {
			c.app.Quit()
		}
		logging.Errorf("AquireData: %v", err)
		dialog.ShowError(err, c.win)
		return
	}
	c.currentPage = c.pages[c.currentPage].Next(c.currentPage)
	c.win.SetContent(c.Window())
	c.pages[c.currentPage].Run()
}

func (c *Wizard) Prev() {
	logging.Debugf("Prev from page %d to %d", c.currentPage, c.pages[c.currentPage].Prev())
	c.currentPage = c.pages[c.currentPage].Prev()
	c.win.SetContent(c.Window())
	c.pages[c.currentPage].Run()
}

func (c *Wizard) Run() {
	c.win.ShowAndRun()
}
