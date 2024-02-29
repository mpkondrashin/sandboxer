package xplatform

import (
	"errors"
	"os"
	"path/filepath"
)

func LinkToStartMenu(folder, name, path string) (string, error) {
	folderPath := filepath.Join(os.Getenv("PROGRAMDATA"), `Microsoft\Windows\Start Menu\Programs`, folder)
	if err := os.Mkdir(folderPath, 0755); err != nil {
		if !errors.Is(err, os.ErrExist) {
			return "", err
		}
	}
	linkPath := filepath.Join(folderPath, name) + ".lnk"
	err := makeLink(path, linkPath)
	if err != nil {
		_ = os.RemoveAll(folderPath)
		return "", err
	}
	return folderPath, nil
}
