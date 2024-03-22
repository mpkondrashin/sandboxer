/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

config_test.go

Small test
*/
package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSave(t *testing.T) {
	c := New("testing_config.yaml")
	//t.Logf("%v", c.SandboxType)
	err := c.Save()
	if err != nil {
		t.Fatal(err)
	}
	err = c.Load()
	if err != nil {
		t.Fatal(err)
	}
}

var configTest = `
sandbox_type: VisionOne
vision_one:
    token: "abc"
    aws_region: "www.com"
analyzer:
    url: ""
    protocol_version: "1.7"
    user_agent: sandboxer/v0.1.5
    product_name: Sandboxer
    hostname: Mikhails-MacBook-Pro.local
    temp_folder: /var/folders/6d/2m86w1vx7hqf411p_05rpq7m0000gn/T/
    source_id: "303"
    source_name: sandboxer
    api_key: ""
    ignore_tls_errors: false
    client_id: ""
proxy:
    type: NTLM
    url: "http://1.1.1.1:8080"
    username: "mike"
    password: "test1234"
    domain: "test.local"
    keepalive: 30s
folder: /Applications
ignore:
    - .DS_Store
    - Thumbs.db
sleep: 5s
periculosum: check
show_password_hint: true
`

func TestLoad(t *testing.T) {
	fileName := "config_test_data.yaml"
	folderName := "testing_load_config"
	if err := os.MkdirAll(folderName, 0755); err != nil {
		t.Fatal(err)
	}
	filePath := filepath.Join(folderName, fileName)
	if err := os.WriteFile(filePath, []byte(configTest), 0644); err != nil {
		t.Fatal(err)
	}
	c := New(filePath)
	if c.Proxy != c.VisionOne.Proxy {
		t.Fatalf("Proxy is %p, but VisionOne.Proxy is %p", c.Proxy, c.VisionOne.Proxy)
	}
	if c.Proxy != c.DDAn.Proxy {
		t.Fatalf("Proxy is %p, but DDAn.Proxy is %p", c.Proxy, c.DDAn.Proxy)
	}
	if err := c.Load(); err != nil {
		t.Fatal(err)
	}
	t.Run("v1 token", func(t *testing.T) {
		actual := c.VisionOne.Token
		expected := "abc"
		if actual != expected {
			t.Errorf("Expected %v, but got %v", expected, actual)
		}
	})
	t.Run("proxy type", func(t *testing.T) {
		actual := c.Proxy.Type
		expected := AuthTypeNTLM
		if actual != expected {
			t.Errorf("Expected %v, but got %v", expected, actual)
		}
	})
	t.Run("v1 proxy type", func(t *testing.T) {
		actual := c.VisionOne.Proxy.Type
		expected := AuthTypeNTLM
		if actual != expected {
			t.Errorf("Expected %v, but got %v", expected, actual)
		}
	})

}

/*
func TestLoad(t *testing.T) {
	conf1 := &Configuration{
		Token:     "1",
		APIKey:    "",
		Region:    "3",
		AccountID: "4",
		//		NSHIRegion:      "5",
	}

	password := "testPasswordtestPassword"
	fileName := "testFile.yaml"

	err := conf1.Save(fileName, password)
	if err != nil {
		t.Fatalf("Failed to load configuration: %s", err)
	}

	conf2 := &Configuration{
		Token:     "1",
		APIKey:    "2",
		Region:    "3",
		AccountID: "4",
		//	NSHIRegion:      "5",
	}
	conf2.Load(fileName, password)

	if conf1.Token != conf2.Token {
		t.Errorf("Expected: %v, got: %v", conf2.Token, conf1.Token)
	}
}

func TestLoadErr(t *testing.T) {
	conf := &Configuration{}
	password := "testPasswordtestPassword"
	fileName := "nonExistentFile.yaml"

	err := conf.Load(fileName, password)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}

	if !os.IsNotExist(err) {
		t.Errorf("Expected 'no such file or directory' error, got: %v", err)
	}
}
*/
