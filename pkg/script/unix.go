package script

import "fmt"

type Unix struct {
}

func (Unix) Extension() string {
	return ".sh"
}

func (Unix) Comment(text string) string {
	return fmt.Sprintf("# %s", text)
}

func (Unix) RemoveDir(path string) string {
	return fmt.Sprintf("rm -r \"%s\"", path)
}
