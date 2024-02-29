package xplatform

import (
	"os"
	"path/filepath"
)

func LinkToStartMenu(folder, name, path string) (string, error) {
	folderPath := filepath.Join(os.Getenv("PROGRAMDATA"), `Microsoft\Windows\Start Menu\Programs`, folder)
	if err := os.Mkdir(folderPath, 0755); err != nil {
		return "", err
	}
	linkPath := filepath.Join(folderPath, name)
	err := makeLink(path, linkPath)
	if err != nil {
		_ = os.RemoveAll(folderPath)
		return "", err
	}
	return folderPath, nil
}
