package main

import (
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

	"sandboxer/pkg/globals"
	"sandboxer/pkg/logging"
)

//go:generate fyne bundle --name IconSVGResource --output resource.go ../../resources/icon.png

type Wizard struct {
	pages          []Page
	current        int
	app            fyne.App
	win            fyne.Window
	installer      *Installer
	capturesFolder string
}

type Page interface {
	Name() string
	Content(win fyne.Window, installer *Installer) fyne.CanvasObject
	Run(win fyne.Window, installer *Installer)
	AquireData(installer *Installer) error
}

func NewWizard(capturesFolder string) *Wizard {
	installer, err := NewInstaller(globals.AppID)
	if err != nil {
		return nil // let it crash
	}
	c := &Wizard{
		app:            app.NewWithID(globals.AppID),
		capturesFolder: capturesFolder,
		installer:      installer,
	}
	c.app.Lifecycle()
	_ = c.installer.LoadConfig()
	c.win = c.app.NewWindow(globals.AppName + " Install Program")
	c.win.Resize(fyne.NewSize(600, 400))
	//c.win.SetFixedSize(true)
	c.win.SetMaster()
	c.pages = []Page{
		&PageIntro{},
		&PageOptions{},
		&PageDomain{},
		&PageFolder{},
		&PageAutostart{},
		&PageInstallation{},
		&PageFinish{},
	}
	prtScr := &desktop.CustomShortcut{KeyName: fyne.KeyI, Modifier: fyne.KeyModifierControl}
	c.win.Canvas().AddShortcut(prtScr, c.captureWindowContents)
	c.win.SetContent(c.Window(c.pages[0]))
	return c
}

func (c *Wizard) captureWindowContents(shortcut fyne.Shortcut) {
	if c.capturesFolder == "" {
		log.Println("--capture is not set")
		return
	}
	fileName := fmt.Sprintf("page_%d.png", c.current)
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
func (c *Wizard) Window(p Page) fyne.CanvasObject {
	left := container.NewVBox()
	image := canvas.NewImageFromResource(ApplicationIcon)
	image.SetMinSize(fyne.NewSize(52, 52))
	image.FillMode = canvas.ImageFillContain
	left.Add(image)
	for _, page := range c.pages {
		if page == p {
			left.Add(widget.NewLabelWithStyle("▶ "+page.Name(), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
		} else {
			left.Add(widget.NewLabel("    " + page.Name()))
		}
	}

	middle := container.NewPadded(container.NewVBox(layout.NewSpacer(), p.Content(c.win, c.installer), layout.NewSpacer()))
	upper := container.NewBorder(nil, nil, container.NewHBox(left, widget.NewSeparator()), nil, middle)
	prevButton, nextButton := c.Buttons()
	buttons := container.NewBorder(nil, nil, nil,
		container.NewHBox(prevButton, nextButton))
	bottom := container.NewVBox(widget.NewSeparator(), buttons)
	_ = bottom

	return container.NewBorder(nil, container.NewPadded(bottom), nil, nil, upper)
}

func (c *Wizard) Buttons() (*widget.Button, *widget.Button) {
	//quitButton := widget.NewButtonWithIcon("Quit", theme.CancelIcon(), c.Quit)

	prevButton := widget.NewButtonWithIcon("Back", theme.NavigateBackIcon(), c.Prev)
	if c.current == 0 {
		prevButton.Disable()
	}

	nextButton := widget.NewButtonWithIcon("Next", theme.NavigateNextIcon(), c.Next)
	nextButton.IconPlacement = widget.ButtonIconTrailingText

	if c.current == len(c.pages)-1 {
		//Button.Disable()
		//nextButton.IconPlacement = widget.ButtonIconTrailingText
		//quitButton.Disable()
		nextButton = widget.NewButtonWithIcon("Quit", theme.CancelIcon(), c.Quit)

	}
	return prevButton, nextButton
}

func (c *Wizard) Quit() {
	logging.Debugf("Quit")
	err := c.pages[c.current].AquireData(c.installer)
	if err != nil {
		logging.Errorf("AquireData: %v", err)
		dialog.ShowError(err, c.win)
	}
	//dialog.ShowConfirm()
	c.app.Quit()
}

func (c *Wizard) Next() {
	logging.Debugf("Next from page %d", c.current)
	err := c.pages[c.current].AquireData(c.installer)
	if err != nil {
		logging.Errorf("AquireData: %v", err)
		dialog.ShowError(err, c.win)
		return
	}
	c.current++
	c.win.SetContent(c.Window(c.pages[c.current]))
	c.pages[c.current].Run(c.win, c.installer)
}

func (c *Wizard) Prev() {
	logging.Debugf("Prev from page %d", c.current)
	c.current--
	c.win.SetContent(c.Window(c.pages[c.current]))
	c.pages[c.current].Run(c.win, c.installer)
}

func (c *Wizard) Run() {
	c.win.ShowAndRun()
}
