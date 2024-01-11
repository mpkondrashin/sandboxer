package script

import "runtime"

type Script interface {
	Extension() string
	Comment(text string) string
	RemoveDir(path string) string
}

func Get() Script {
	if runtime.GOOS == "windows" {
		return Windows{}
	}
	return Unix{}
}
