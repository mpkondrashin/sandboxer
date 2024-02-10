package logging

import (
	"fmt"
	"os"
	"path/filepath"

	"testing"
)

func TestRotatedFile_Rotate(t *testing.T) {
	dir := "testing_rotated"
	if err := os.RemoveAll(dir); err != nil {
		t.Fatal(t)
	}
	if err := os.Mkdir(dir, 0775); err != nil {
		t.Fatal(t)
	}

	fileName := "some.log"
	r := NewRotated(dir, fileName, 0666, 10, 2)
	if err := r.Open(); err != nil {
		t.Fatal(t)
	}
	defer r.Close()
	logIt := func(data string) {
		n, err := fmt.Fprint(r, data)
		if err != nil {
			t.Errorf("Got %v", err)
		} else if n != len(data) {
			t.Errorf("Expected %v but got %v", len(data), n)
		}
	}
	logIt("0123")
	logIt("abcd")
	logIt("ABCD")
	logIt("0123456789")
	check := func(fileName, data string) {
		filePath := filepath.Join(dir, fileName)
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Errorf("%v: %v", fileName, err)
			return
		}
		if data != string(content) {
			t.Errorf("Expected %v but got %v", data, string(content))
		}
	}
	check("some.log", "0123456789")
	check("some.log.0", "ABCD")
	check("some.log.1", "0123abcd")
	logIt("9876543210")
	check("some.log", "9876543210")
	check("some.log.0", "0123456789")
	check("some.log.1", "ABCD")

	filePath := filepath.Join(dir, "some.log.2")
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Errorf("%v exists", filePath)
	}
}
