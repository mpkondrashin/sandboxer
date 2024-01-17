package script

import "runtime"

type Platform interface {
	Extension() string
	Comment(text string) string
	RemoveDir(path string) string
	UninstallService(name string) string
	StopService(name string) string
}

func Get() Platform {
	if runtime.GOOS == "windows" {
		return Windows{}
	}
	return Darwin{}
}
