/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
This software is distributed under MIT license as stated in LICENSE file

page_finish.go

Final installer page
*/
package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"sandboxer/pkg/globals"
	"sandboxer/pkg/xplatform"
)

type PageFinish struct {
	BasePage
	//runCheck *widget.Check
}

var _ Page = &PageFinish{}

func (p *PageFinish) Name() string {
	return "Finish"
}

func (p *PageFinish) Next(previousPage PageIndex) PageIndex {
	p.SavePrevious(previousPage)
	return pgExit
}

func (p *PageFinish) Content() fyne.CanvasObject {
	l1 := widget.NewLabel(globals.AppName + " service sucessfully installed.")
	hint := "Right click on any file and pick Send To -> " + globals.AppName + "."
	if !xplatform.IsWindows() {
		hint = "Right click on any file and pick Quick Actions -> " + globals.AppName + "."
	}
	l2 := widget.NewLabel(hint)
	//p.runCheck = widget.NewCheck("Run "+globals.AppName+" now", nil)
	//p.runCheck.SetChecked(true)
	return container.NewVBox(l1, l2) //, p.runCheck)
}

//func (p *PageFinish) Run() {}

//func (p *PageFinish) AquireData(installer *Installer) error {
//if !p.runCheck.Checked {
//	return nil
//}
//path := filepath.Join(installer.InstallFolder(), globals.Name+".exe")
//cmd := exec.Command(path, "--submissions")
//return cmd.Start()
//	return nil
//}
