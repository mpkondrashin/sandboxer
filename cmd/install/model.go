package main

import (
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"

	"examen/pkg/config"
)

type Model struct {
	appID    string
	fileName string
	config   config.Configuration
	hash     string
}

func NewModel(appID string) *Model {
	folder, _ := config.InstallFolder()
	return &Model{appID: appID,
		config: config.Configuration{Folder: folder}}
}

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

func (m *Model) LoadConfig() error {
	filePath, err := m.configFilePath()
	if err != nil {
		return err
	}
	err = m.config.Load(filePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (m *Model) CalculateHash() string {
	h := sha1.New()
	h.Write([]byte(m.config.Token))
	h.Write([]byte(m.config.Domain))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (m *Model) Load() error {
	filePath, err := m.configFilePath()
	if err != nil {
		return err
	}

	if err := m.config.Load(filePath); err != nil {
		return err
	}
	m.hash = m.CalculateHash()
	return nil
}

func (m *Model) Save() error {
	filePath, err := m.configFilePath()
	if err != nil {
		return err
	}
	if err := m.config.Save(filePath); err != nil {
		return err
	}
	m.hash = m.CalculateHash()
	return nil
}

func (m *Model) Changed() bool {
	return m.hash != m.CalculateHash()
}

func (m *Model) configFilePath() (string, error) {
	folder, err := m.config.ConfigFileFolder(m.appID)
	if err != nil {
		return "", err
	}
	return filepath.Join(folder, m.fileName), nil
}
