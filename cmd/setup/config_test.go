package main

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	conf1 := &Configuration{
		token:     "1",
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
		token:     "1",
		APIKey:    "2",
		Region:    "3",
		AccountID: "4",
		//	NSHIRegion:      "5",
	}
	conf2.Load(fileName, password)

	if conf1.token != conf2.token {
		t.Errorf("Expected: %v, got: %v", conf2.token, conf1.token)
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
