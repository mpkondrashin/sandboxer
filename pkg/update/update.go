package update

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

func DownloadRelease(version, filename, folder string, progress func(float32) error) error {
	url := fmt.Sprintf("https://github.com/mpkondrashin/sandboxer/releases/download/%s/%s", version, filename)
	client := http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return err
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
