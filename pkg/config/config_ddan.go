package config

import (
	"errors"
	"net/url"
	"os"
	"sandboxer/pkg/globals"
	"sync"

	"github.com/google/uuid"
	"github.com/mpkondrashin/ddan"
)

type DDAn struct {
	mx              sync.RWMutex `gsetter:"-"`
	URL             string       `yaml:"url"`
	ProtocolVersion string       `yaml:"protocol_version"`
	UserAgent       string       `yaml:"user_agent"`
	ProductName     string       `yaml:"product_name"`
	Hostname        string       `yaml:"hostname"`
	TempFolder      string       `yaml:"temp_folder"`
	SourceID        string       `yaml:"source_id"`
	SourceName      string       `yaml:"source_name"`
	APIKey          string       `yaml:"api_key"`
	IgnoreTLSErrors bool         `yaml:"ignore_tls_errors"`
	ClientUUID      string       `yaml:"-"`
	Proxy           *Proxy       `yaml:"-" gsetter:"-"`
}

func NewDefaultDDAn(proxy *Proxy) *DDAn {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = err.Error()
	}
	return &DDAn{
		ProtocolVersion: "1.8",
		UserAgent:       globals.Name + "/" + globals.Version,
		ProductName:     globals.AppName,
		Hostname:        hostname,
		TempFolder:      os.TempDir(),
		SourceID:        "303",
		SourceName:      globals.Name,
		IgnoreTLSErrors: false,
		Proxy:           proxy,
	}
}

func (d *DDAn) Update(newDDAn *DDAn) {
	d.mx.Lock()
	defer d.mx.Unlock()
	newDDAn.mx.RLock()
	defer newDDAn.mx.RUnlock()
	d.URL = newDDAn.URL
	d.APIKey = newDDAn.APIKey
	d.IgnoreTLSErrors = newDDAn.IgnoreTLSErrors
}

func (d *DDAn) Analyzer() (*ddan.Client, error) {
	d.mx.RLock()
	defer d.mx.RUnlock()
	u, err := url.Parse(d.URL)
	if err != nil {
		return nil, err
	}
	if d.Hostname == "" {
		d.Hostname, err = os.Hostname()
		if err != nil {
			return nil, err
		}
	}
	analyzer := ddan.NewClient(d.ProductName, d.Hostname)
	analyzer.SetAnalyzer(u, d.APIKey, d.IgnoreTLSErrors)
	analyzer.SetSource(d.SourceID, d.SourceName)
	analyzer.SetUUID(d.ClientUUID)
	analyzer.SetProtocolVersion(d.ProtocolVersion)

	if d.Proxy == nil {
		return analyzer, nil
	}
	modifier, err := d.Proxy.Modifier()
	if err != nil {
		return nil, err
	}
	analyzer.ModifyTransport(modifier)

	return analyzer, nil
}

func (d *DDAn) AnalyzerWithUUID() (*ddan.Client, error) {
	if err := d.ProvideUUID(); err != nil {
		return nil, err
	}
	return d.Analyzer()
}

func (d *DDAn) ProvideUUID() (err error) {
	d.mx.Lock()
	defer d.mx.Unlock()
	if d.ClientUUID == "" {
		d.ClientUUID, err = GenerateUUID()
	}
	return
}

func (d *DDAn) LoadClientUUID() (err error) {
	path, err := globals.AnalyzerClientUUIDFilePath()
	if err != nil {
		return err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	d.ClientUUID = string(data)
	return nil
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
