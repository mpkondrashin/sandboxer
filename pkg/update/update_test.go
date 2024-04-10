/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

update_test.go

Test updates
*/
package update

import (
	"errors"
	"testing"
)

const repo = "sandboxer"

func TestCheckUpdate(t *testing.T) {
	u, err := CheckLocationGithub(repo)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("url = %v", u)
	version := ParseVersion(u)
	t.Logf("version = %v", version)
	err = DownloadRelease(version, "Setup.zip", ".", func(p float32) error {
		t.Logf("Downloaded %d%%", int(p*100))
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}

func TestCheckRelease(t *testing.T) {
	version := "v100.100.100"
	filename := "setup_darwin_arm64.zip"
	err := CheckRelease(version, filename)
	t.Logf("error: %v", err)
	if err == nil {
		t.Errorf("succeeded to download version %s", version)
	}
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("wrong error type: %v", err)
	}
	version = "v0.4.4"
	err = CheckRelease(version, filename)
	t.Logf("error: %v", err)
	if err != nil {
		t.Errorf("error checking existing version %s", version)
	}

}
