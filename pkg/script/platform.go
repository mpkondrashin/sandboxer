/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

platform.go

uninstall script interface
*/
package script

import "runtime"

type Platform interface {
	Extension() string
	Comment(text string) string
	RemoveDir(path string) string
	UninstallService(name string) string
	StopService(name string) string
	StopProcess(name string) string
}

func Get() Platform {
	if runtime.GOOS == "windows" {
		return Windows{}
	}
	return Darwin{}
}
