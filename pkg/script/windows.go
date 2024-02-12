/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

windows.go

uninstall script commands for Windows
*/
package script

import "fmt"

type Windows struct {
}

func (Windows) Extension() string {
	return ".cmd"
}

func (Windows) Comment(text string) string {
	return fmt.Sprintf("rem %s", text)
}

func (Windows) RemoveDir(path string) string {
	return fmt.Sprintf("del /F /S /Q \"%s\"", path)
}

func (Windows) UninstallService(name string) string {
	return fmt.Sprintf("sc delete %s", name)
}

func (Windows) StopService(name string) string {
	return fmt.Sprintf("sc stop %s", name)
}

func (Windows) StopProcess(name string) string {
	return fmt.Sprintf("taskkill /im %s.exe /f", name)
}
