//go:build darwin

/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

alert_darwin.go

Show alert on macOS
*/

package xplatform

import (
	"fmt"
	"os/exec"
)

func Alert(title, subtitle, message, _ string) error {
	osa, err := exec.LookPath("osascript")
	if err != nil {
		return err
	}
	script := fmt.Sprintf(`display notification "%s" with title "%s" subtitle "%s" sound name "default" with icon`, message, title, subtitle)
	return exec.Command(osa, "-e", script).Run()
}
