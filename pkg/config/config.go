package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Token  string `yaml:"token"`
	Domain string `yaml:"aws_region"`
	Folder string `yaml:"folder"`
}

// Save - writes Configuration struct to file as YAML
func (c *Configuration) Save(fileName string) (err error) {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, data, 0600)
}

// Load - reads Configuration struct from YAML file
func (c *Configuration) Load(fileName string) error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, c)
}

var (
	ErrNoUserProfile = errors.New("missing environment variable")
	ErrUnsupportedOS = errors.New("unsupported OS")
)

func (c *Configuration) ConfigFileFolder(appID string) (string, error) {
	if runtime.GOOS == "windows" {
		return c.configFileFolder("USERPROFILE", "AppData\\Local", appID)
	}
	if runtime.GOOS == "linux" {
		return c.configFileFolder("HOME", ".config", appID)
	}
	if runtime.GOOS == "darwin" {
		return c.configFileFolder("HOME", "Library/Application Support", appID)
	}
	return "", fmt.Errorf("%s: %w", runtime.GOOS, ErrUnsupportedOS)
}

func (c *Configuration) configFileFolder(profileVariable string, dir string, appID string) (string, error) {
	userProfile := os.Getenv(profileVariable)
	if userProfile == "" {
		return "", fmt.Errorf("%s: %w", profileVariable, ErrNoUserProfile)
	}
	folder := filepath.Join(userProfile, dir, appID)
	err := os.MkdirAll(folder, 0700)
	if err != nil {
		return "", err
	}
	return folder, nil
}

func InstallFolder() (string, error) {
	if runtime.GOOS == "windows" {
		return os.Getenv("PROGRAMFILES"), nil
	}
	if runtime.GOOS == "linux" {
		return "/usr/local/bin", nil
	}
	if runtime.GOOS == "darwin" {
		return "/Applications", nil
	}
	return "", fmt.Errorf("%s: %w", runtime.GOOS, ErrUnsupportedOS)
}
