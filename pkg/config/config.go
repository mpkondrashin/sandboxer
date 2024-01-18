package config

import (
	"errors"
	"examen/pkg/globals"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Token       string        `yaml:"token"`
	Domain      string        `yaml:"aws_region"`
	Folder      string        `yaml:"folder"`
	Ignore      []string      `yaml:"ignore"`
	Periculosum string        `yaml:"periculosum"`
	Sleep       time.Duration `yaml:"sleep"`
}

func New() *Configuration {
	return &Configuration{
		Folder:      InstallFolder(),
		Ignore:      []string{".DS_Store", "Thumbs.db"},
		Periculosum: "check",
		Sleep:       5 * time.Second,
	}
}

func (c *Configuration) LogFolder() string {
	return filepath.Join(c.Folder, globals.AppFolderName, "logs")
}

func LoadConfiguration(appID string, fileName string) (*Configuration, error) {
	folder, err := ConfigFileFolder(appID)
	if err != nil {
		return nil, err
	}
	filePath := filepath.Join(folder, fileName)
	c := &Configuration{}
	if err := c.Load(filePath); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Configuration) PericulosumPath() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(path), c.Periculosum), nil
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

func ConfigFileFolder(appID string) (string, error) {
	if runtime.GOOS == "windows" {
		return configFileFolder("USERPROFILE", "AppData\\Local", appID)
	}
	if runtime.GOOS == "linux" {
		return configFileFolder("HOME", ".config", appID)
	}
	if runtime.GOOS == "darwin" {
		return configFileFolder("HOME", "Library/Application Support", appID)
	}
	return "", fmt.Errorf("%s: %w", runtime.GOOS, ErrUnsupportedOS)
}

func configFileFolder(profileVariable string, dir string, appID string) (string, error) {
	userProfile := os.Getenv(profileVariable)
	if userProfile == "" {
		return "", fmt.Errorf("%s: %w", profileVariable, ErrNoUserProfile)
	}
	folder := filepath.Join(userProfile, dir, appID)
	//err := os.MkdirAll(folder, 0700)
	//if err != nil {
	//return "", err
	//}
	return folder, nil
}

func InstallFolder() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("PROGRAMFILES")
	}
	if runtime.GOOS == "linux" {
		return "/usr/local/bin"
	}
	if runtime.GOOS == "darwin" {
		return "/Applications"
	}
	return ""
}
