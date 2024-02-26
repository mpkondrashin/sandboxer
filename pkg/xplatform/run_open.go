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
