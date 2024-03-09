/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

run_open.go

Run open/start commands
*/

package xplatform

import (
	"os/exec"
	"runtime"
)

func RunOpen(path string) error {
	name := "open"
	args := []string{path}
	if runtime.GOOS == "windows" {
		name = "cmd"
		args = []string{"/C", "start", path}
	}
	cmd := exec.Command(name, args...)
	return cmd.Run()
}
