package main

import (
	"os"
	"path/filepath"

	"examen/pkg/config"
)

type Installer struct {
	appID    string
	fileName string
	config   config.Configuration
	//hash     string
}

func NewInstaller(appID string) *Installer {
	folder, _ := config.InstallFolder()
	return &Installer{appID: appID,
		config: config.Configuration{Folder: folder}}
}

/*
	func (m *Model) ConfigExists() (bool, error) {
		path, err := m.configFilePath()
		if err != nil {
			return false, err
		}
		_, err = os.Stat(path)
		if err == nil {
			return true, nil
		}
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
*/
func (i *Installer) LoadConfig() error {
	filePath, err := i.configFilePath()
	if err != nil {
		return err
	}
	err = i.config.Load(filePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

/*
func (m *Model) CalculateHash() string {
	h := sha1.New()
	h.Write([]byte(m.config.Token))
	h.Write([]byte(m.config.Domain))
	return fmt.Sprintf("%x", h.Sum(nil))
}
*/

func (i *Installer) Load() error {
	filePath, err := i.configFilePath()
	if err != nil {
		return err
	}

	if err := i.config.Load(filePath); err != nil {
		return err
	}
	//m.hash = m.CalculateHash()
	return nil
}

func (i *Installer) Save() error {
	filePath, err := i.configFilePath()
	if err != nil {
		return err
	}
	if err := i.config.Save(filePath); err != nil {
		return err
	}
	//m.hash = m.CalculateHash()
	return nil
}

/*
	func (m *Model) Changed() bool {
		return m.hash != m.CalculateHash()
	}
*/
func (i *Installer) configFilePath() (string, error) {
	folder, err := i.config.ConfigFileFolder(i.appID)
	if err != nil {
		return "", err
	}
	return filepath.Join(folder, i.fileName), nil
}
