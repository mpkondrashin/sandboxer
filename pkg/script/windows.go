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
