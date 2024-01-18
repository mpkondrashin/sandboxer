package script

import "fmt"

type Darwin struct {
}

func (Darwin) Extension() string {
	return ".sh"
}

func (Darwin) Comment(text string) string {
	return fmt.Sprintf("# %s", text)
}

func (Darwin) RemoveDir(path string) string {
	return fmt.Sprintf("rm -r \"%s\"", path)
}

func (Darwin) UninstallService(name string) string {
	return fmt.Sprintf("launchctl unload /System/Library/LaunchDaemons/%s.plist", name)
}

func (Darwin) StopService(name string) string {
	return fmt.Sprintf("# stopping %s", name)
}
