package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/kardianos/service"
	"gopkg.in/yaml.v2"

	"examen/pkg/globals"
)

type VisionOne struct {
	Token  string        `yaml:"token"`
	Domain string        `yaml:"aws_region"`
	Sleep  time.Duration `yaml:"sleep"`
}

type Configuration struct {
	filePath    string
	VisionOne   VisionOne `yaml:"vision_one"`
	Folder      string    `yaml:"folder"`
	Ignore      []string  `yaml:"ignore"`
	Periculosum string    `yaml:"periculosum"`
}

func New(filePath string) *Configuration {
	return &Configuration{
		filePath:    filePath,
		Folder:      globals.InstallFolder(),
		Ignore:      []string{".DS_Store", "Thumbs.db"},
		Periculosum: "check",
		VisionOne: VisionOne{
			Sleep: 5 * time.Second,
		},
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

func (c *Configuration) Path(fileName string) string {
	return filepath.Join(c.Folder, globals.AppFolderName, fileName)
}

func (c *Configuration) Service(i service.Interface) (service.Service, error) {
	svcConfig := &service.Config{
		Name:        globals.SvcName,
		DisplayName: globals.SvcDisplayName,
		Description: globals.SvcDescription,
		Executable:  c.Path(globals.SvcFileName),
	}
	return service.New(i, svcConfig)
}
