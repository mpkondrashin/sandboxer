/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

task_test.go

Test the basic task functions
*/
package task

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHashes(t *testing.T) {
	text := "Hello World!"
	folder := filepath.Join("tests", "hashes")
	if err := os.MkdirAll(folder, 0755); err != nil {
		t.Fatal(err)
	}
	fileName := "file.txt"
	filePath := filepath.Join(folder, fileName)
	if err := os.WriteFile(filePath, []byte(text), 0644); err != nil {
		t.Fatal(err)
	}
	tsk := NewTask(0, FileTask, filePath)
	if err := tsk.CalculateHash(); err != nil {
		t.Fatal(err)
	}

	var actual, expected string
	expected = "7f83b1657ff1fc53b92dc18148a1d65dfc2d4b1fa3d677284addd200126d9069"
	actual = tsk.SHA256
	if expected != actual {
		t.Errorf("Expected %s, but got %s", expected, actual)
	}

	expected = "2ef7bde608ce5404e97d5f042f95f89f1c232871"
	actual = tsk.SHA1
	if expected != actual {
		t.Errorf("Expected %s, but got %s", expected, actual)
	}

	expected = "ed076287532e86365e841e92bfc50d8c"
	actual = tsk.MD5
	if expected != actual {
		t.Errorf("Expected %s, but got %s", expected, actual)
	}
}
