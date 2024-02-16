package update

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sandboxer/pkg/globals"
	"sandboxer/pkg/logging"
	"strings"
	"time"

	"golang.org/x/mod/semver"
)

const (
// repo = "sandboxer"
)

func CheckLocationGithub(repo string) (string, error) {
	checkURL := fmt.Sprintf("https://github.com/mpkondrashin/%s/releases/latest", repo)
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Get(checkURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return resp.Header.Get("Location"), nil
}

func ParseVersion(releaseURL string) string {
	// https://github.com/mpkondrashin/sandboxer/releases/tag/v0.0.10
	index := strings.LastIndex(releaseURL, "/")
	if index == -1 {
		return ""
	}
	return releaseURL[index+1:]
}

var ErrMissingVersion = errors.New("missing version")

func LatestVersion(repo string) (string, error) {
	loc, err := CheckLocationGithub(repo)
	if err != nil {
		return "", err
	}
	ver := ParseVersion(loc)
	if ver == "" {
		return "", fmt.Errorf("%s: %w", loc, ErrMissingVersion)
	}
	return ver, nil
}

var ErrNotFound = errors.New("not found")

func DownloadRelease(version, filename, folder string, progress func(float32) error) error {
	url := fmt.Sprintf("https://github.com/mpkondrashin/sandboxer/releases/download/%s/%s", version, filename)
	client := http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("%s: %w", url, ErrNotFound)
	}
	defer resp.Body.Close()
	filePath := filepath.Join(folder, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	var progressInt func(c int64) error
	if progress != nil {
		progressInt = func(c int64) error {
			return progress(float32(c) / float32(resp.ContentLength))
		}
	}
	return Download(file, resp.Body, progressInt)
}

var ErrInvalidWrite = errors.New("invalid write result")

func Download(dst io.Writer, src io.Reader, progress func(int64) error) (err error) {
	size := 32 * 1024
	buf := make([]byte, size)
	var written int64
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = ErrInvalidWrite
				}
			}
			written += int64(nw)
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if progress != nil {
			progress(written)
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return err
}

type Update struct {
	Version string
	Date    time.Time
}

func NeedUpdateWindow() (bool, error) {
	version, err := LatestVersion(globals.Name)
	logging.Debugf("NeedUpdateWindow:LatestVersion(%s): %s, %v", globals.Name, version, err)
	if err != nil {
		return false, err
	}
	cmp := semver.Compare(version, globals.Version)
	logging.Debugf("NeedUpdateWindow:semver.Compare(%s, %s): %v", version, globals.Version, cmp)
	//logging.Debugf("Compare %s vs %s: %d", version, globals.Version, cmp)
	if cmp != 1 {
		return false, nil
	}
	folder, err := globals.UserDataFolder()
	logging.Debugf("NeedUpdateWindow:globals.UserDataFolder(): %s %v", folder, err)
	if err != nil {
		return false, err
	}
	fileName := "check_update.json"
	filePath := filepath.Join(folder, fileName)
	logging.Debugf("NeedUpdateWindow:filePath: %s", filePath)
	data, err := os.ReadFile(filePath)
	logging.Debugf("NeedUpdateWindow:ReadFile(%s): %s, %v", filePath, string(data), err)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			u := Update{
				Version: version,
				Date:    time.Now(),
			}
			data, err = json.MarshalIndent(u, "", "    ")
			er := os.WriteFile(filePath, data, 0644)
			logging.Debugf("NeedUpdateWindow:WriteFile(%s, %s): %v", filePath, string(data), er)
			if er != nil {
				return false, er
			}
			return true, err
		}
		return false, err
	}
	var u Update
	err = json.Unmarshal(data, &u)
	if err != nil {
		return false, err
	}
	cmp = semver.Compare(version, u.Version)
	logging.Debugf("NeedUpdateWindow:Compare(%s, %s): %d", version, u.Version, cmp)
	if cmp == 0 {
		return false, nil
	}
	u.Version = version
	data, err = json.MarshalIndent(u, "", "    ")
	if err != nil {
		return false, err
	}
	err = os.WriteFile(fileName, data, 0644)
	logging.Debugf("NeedUpdateWindow:WriteFile(%s, %s): %d", fileName, data, err)
	if err != nil {
		return false, err
	}
	return true, nil
}
