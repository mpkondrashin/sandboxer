package main

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type Model struct {
	appName  string
	fileName string
	password string
	config   Configuration
	hash     string
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
	err = m.config.Load(filePath, m.password)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (m *Model) CalculateHash() string {
	h := sha1.New()
	h.Write([]byte(m.password))
	h.Write([]byte(m.config.APIKey))
	h.Write([]byte(m.config.Region))
	h.Write([]byte(m.config.AccountID))
	h.Write([]byte(m.config.Domain))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (m *Model) Load() error {
	filePath, err := m.configFilePath()
	if err != nil {
		return err
	}

	if err := m.config.Load(filePath, m.password); err != nil {
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
	if err := m.config.Save(filePath, m.password); err != nil {
		return err
	}
	m.hash = m.CalculateHash()
	return nil
}

func (m *Model) Changed() bool {
	return m.hash != m.CalculateHash()
}

func (m *Model) configFilePath() (string, error) {
	folder, err := m.ConfigFileFolder()
	if err != nil {
		return "", err
	}
	return filepath.Join(folder, m.fileName), nil
}

var (
	ErrNoUserProfile = errors.New("missing environment variable")
	ErrUnsupportedOS = errors.New("unsupported OS")
)

func (m *Model) ConfigFileFolder() (string, error) {
	if runtime.GOOS == "windows" {
		return m.configFileFolder("USERPROFILE", "AppData\\Local")
	}
	if runtime.GOOS == "linux" {
		return m.configFileFolder("HOME", ".config")
	}
	if runtime.GOOS == "darwin" {
		return m.configFileFolder("HOME", "Library/Application Support")
	}
	return "", fmt.Errorf("%s: %w", runtime.GOOS, ErrUnsupportedOS)
}

func (m *Model) configFileFolder(profileVariable string, dir string) (string, error) {
	userProfile := os.Getenv(profileVariable)
	if userProfile == "" {
		return "", fmt.Errorf("%s: %w", profileVariable, ErrNoUserProfile)
	}
	folder := filepath.Join(userProfile, dir, m.appName)
	err := os.MkdirAll(folder, 0700)
	if err != nil {
		return "", err
	}
	return folder, nil
}
