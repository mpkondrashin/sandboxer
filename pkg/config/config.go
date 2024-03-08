/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

config.go

Configuration
*/
package config

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/mpkondrashin/ddan"
	"github.com/mpkondrashin/vone"
	"gopkg.in/yaml.v3"

	"sandboxer/pkg/globals"
	"sandboxer/pkg/xplatform"
)

type VisionOne struct {
	Token  string `yaml:"token"`
	Domain string `yaml:"domain"`
}

func (s *VisionOne) VisionOneSandbox() (*vone.VOne, error) {
	token := s.Token
	if token == "" {
		return nil, errors.New("token is not set")
	}
	domain := s.Domain
	if domain == "" {
		return nil, errors.New("domain is not set")
	}
	return vone.NewVOne(domain, token), nil
}

type DDAn struct {
	URL             string `yaml:"url"`
	ProtocolVersion string `yaml:"protocol_version"`
	UserAgent       string `yaml:"user_agent"`
	ProductName     string `yaml:"product_name"`
	Hostname        string `yaml:"hostname"`
	TempFolder      string `yaml:"temp_folder"`
	SourceID        string `yaml:"source_id"`
	SourceName      string `yaml:"source_name"`
	APIKey          string `yaml:"api_key"`
	IgnoreTLSErrors bool   `yaml:"ignore_tls_errors"`
	ClientUUID      string `yaml:"-"`
}

func NewDefaultDDAn() DDAn {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = err.Error()
	}
	return DDAn{
		ProtocolVersion: "1.8",
		UserAgent:       globals.Name + "/" + globals.Version,
		ProductName:     globals.AppName,
		Hostname:        hostname,
		TempFolder:      os.TempDir(),
		SourceID:        "303",
		SourceName:      globals.Name,
		//		APIKey          string `yaml:"api_key"`
		IgnoreTLSErrors: false,
		//ClientUUID      string `yaml:"client_uuid"`
	}
}

func (d *DDAn) Analyzer() (*ddan.Client, error) {
	u, err := url.Parse(d.URL)
	if err != nil {
		return nil, err
	}
	analyzer := ddan.NewClient(d.ProductName, d.Hostname)
	analyzer.SetAnalyzer(u, d.APIKey, d.IgnoreTLSErrors)
	analyzer.SetSource(d.SourceID, d.SourceName)
	analyzer.SetUUID(d.ClientUUID)
	analyzer.SetProtocolVersion(d.ProtocolVersion)
	return analyzer, nil
}

func (d *DDAn) AnalyzerWithUUID() (*ddan.Client, error) {
	var err error
	if d.ClientUUID == "" {
		if d.ClientUUID == "" {
			d.ClientUUID, err = GenerateUUID()
			if err != nil {
				return nil, err
			}
		}
	}
	return d.Analyzer()
}

func GenerateUUID() (string, error) {
	path, err := globals.AnalyzerClientUUIDFilePath()
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		//logging.Errorf("Client UUID file read error: %v", err)
		if !errors.Is(err, os.ErrNotExist) {
			return "", err
		}
		uuidData, err := uuid.NewRandom()
		if err != nil {
			return "", err
		}
		err = os.WriteFile(path, []byte(uuidData.String()), 0600)
		if err != nil {
			return "", err
		}
		return uuidData.String(), nil
	}
	return string(data), nil
}

type Configuration struct {
	filePath         string
	Version          string
	SandboxType      SandboxType   `yaml:"sandbox_type"`
	VisionOne        VisionOne     `yaml:"vision_one"`
	DDAn             DDAn          `yaml:"analyzer"`
	Folder           string        `yaml:"folder"`
	Ignore           []string      `yaml:"ignore"`
	Sleep            time.Duration `yaml:"sleep"`
	Periculosum      string        `yaml:"periculosum"`
	ShowPasswordHint bool          `yaml:"show_password_hint"`
	TasksKeepDays    int           `yaml:"task_keep_days"`
}

func New(filePath string) *Configuration {
	return &Configuration{
		filePath:         filePath,
		Version:          "",
		SandboxType:      SandboxVisionOne,
		Folder:           xplatform.InstallFolder(),
		Ignore:           []string{".DS_Store", "Thumbs.db"},
		Periculosum:      "check",
		ShowPasswordHint: true,
		TasksKeepDays:    60,
		Sleep:            5 * time.Second,
		VisionOne:        VisionOne{},
		DDAn:             NewDefaultDDAn(),
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
