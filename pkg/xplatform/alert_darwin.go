//go:build darwin

package xplatform

import (
	"fmt"
	"os/exec"
)

func Alert(_, title, subtitle, message string) error {
	osa, err := exec.LookPath("osascript")
	if err != nil {
		return err
	}
	script := fmt.Sprintf(`display notification "%s" with title "%s" subtitle "%s" sound name "default" with icon`, message, title, subtitle)
	return exec.Command(osa, "-e", script).Run()
}
