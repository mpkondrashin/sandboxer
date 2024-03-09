/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

start_menu.go

Add link to Windows start menu
*/
package xplatform

import (
	"errors"
	"os"
	"path/filepath"
)

func LinkToStartMenu(dryRun bool, folder, name, path string, asAdministrator bool) (string, error) {
	folderPath := filepath.Join(os.Getenv("PROGRAMDATA"), `Microsoft\Windows\Start Menu\Programs`, folder)
	linkPath := filepath.Join(folderPath, name) + ".lnk"
	if dryRun {
		return linkPath, nil
	}
	if err := os.Mkdir(folderPath, 0755); err != nil {
		if !errors.Is(err, os.ErrExist) {
			return "", err
		}
	}
	err := makeLink(path, linkPath, asAdministrator)
	if err != nil {
		_ = os.RemoveAll(folderPath)
		return "", err
	}
	return folderPath, nil
}
