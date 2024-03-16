/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

config.go

Configuration
*/
package config

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"gopkg.in/yaml.v3"

	"sandboxer/pkg/globals"
	"sandboxer/pkg/xplatform"
)

type Configuration struct {
	mx                sync.RWMutex `gsetter:"-"`
	filePath          string
	Version           string
	SandboxType       SandboxType   `yaml:"sandbox_type"`
	VisionOne         *VisionOne    `yaml:"vision_one" gsetter:"-"`
	DDAn              *DDAn         `yaml:"analyzer" gsetter:"-"`
	Folder            string        `yaml:"folder"`
	Ignore            []string      `yaml:"ignore"`
	Sleep             time.Duration `yaml:"sleep"`
	Periculosum       string        `yaml:"periculosum"`
	ShowPasswordHint  bool          `yaml:"show_password_hint"`
	TasksKeepDays     int           `yaml:"task_keep_days"`
	ShowNotifications bool          `yaml:"notifications"`
}

func New(filePath string) *Configuration {
	return &Configuration{
		filePath:          filePath,
		Version:           "",
		SandboxType:       SandboxVisionOne,
		Folder:            xplatform.InstallFolder(),
		Ignore:            []string{".DS_Store", "Thumbs.db"},
		Periculosum:       "check",
		ShowPasswordHint:  true,
		TasksKeepDays:     60,
		Sleep:             5 * time.Second,
		VisionOne:         &VisionOne{},
		DDAn:              NewDefaultDDAn(),
		ShowNotifications: true,
	}
}

func (c *Configuration) GetFilePath() string {
	return c.filePath
}

/*
func (c *Configuration) Load() error {
	folder, err := ConfigFileFolder(appID)
	if err != nil {
		return nil, err
	}
	c := &Configuration{
		filePath: filepath.Join(folder, fileName),
	}
	if err := c.Load(); err != nil {
		return nil, err
	}
	return c, nil
}*/

func (c *Configuration) PericulosumPath() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(path), c.Periculosum), nil
}

// Save - writes Configuration struct to file as YAML
func (c *Configuration) Save() (err error) {
	c.Version = globals.Version
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(c.filePath, data, 0600)
}

// Load - reads Configuration struct from YAML file
func (c *Configuration) Load() error {
	data, err := os.ReadFile(c.filePath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, c)
}

/*
	func (c *Configuration) Path(fileName string) string {
		return filepath.Join(c.Folder, globals.AppFolderName, fileName)
	}
*/
func (c *Configuration) Resource(fileName string) string {
	if xplatform.IsWindows() {
		return filepath.Join(c.Folder, globals.AppFolderName, fileName)
	} else {
		return filepath.Join(c.Folder, globals.AppFolderName+".app", "Contents", "Resources", fileName)
	}
}

/*
func (c *Configuration) Service(i service.Interface) (service.Service, error) {
	svcConfig := &service.Config{
		Name:        globals.SvcName,
		DisplayName: globals.SvcDisplayName,
		Description: globals.SvcDescription,
		Executable:  c.Path(globals.SvcFileName),
	}
	return service.New(i, svcConfig)
}
*/
