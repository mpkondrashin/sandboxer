/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

darwin.go

uninstall script commands for MacOS
*/
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

func (Darwin) StopProcess(name string) string {
	return fmt.Sprintf("killall %s", name)
}
