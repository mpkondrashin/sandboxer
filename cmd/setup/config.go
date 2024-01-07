package main

import (
	"errors"
	"os"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	token     string
	APIKey    string `yaml:"api_key"`
	Region    string `yaml:"region"`
	AccountID string `yaml:"account_id"`
	Domain    string `yaml:"aws_region"`
}

// Save - writes Configuration struct to file as YAML
func (c *Configuration) Save(fileName, password string) (err error) {
	//log.Print("XXX Save" + fileName + password + c.AccountID)
	c.APIKey, err = ReliableEncrypt(c.token, password)
	if err != nil {
		return err
	}
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, data, 0600)
}

// Load - reads Configuration struct from YAML file
func (c *Configuration) Load(fileName, password string) error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	//	log.Println("data:", string(data))
	//	log.Println("API:", c.APIKey)
	if err := yaml.Unmarshal(data, c); err != nil {
		return err
	}
	c.token, err = ReliableDecrypt(c.APIKey, password)
	//log.Println("decrypted:", c.apiKeyDecrypted, "err", err)
	return err //if err != nil {
	//return err
	//}
}

func ValidateAccountID(s string) error {
	if len(s) != 12 {
		return errors.New("Account ID must be 12 digits")
	}
	for _, char := range s {
		if char < '0' || char > '9' {
			return errors.New("Account ID must only contain digits")
		}
	}
	return nil
}
