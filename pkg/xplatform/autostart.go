package xplatform

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"strings"
)

func AutoStart(name string, appPath string) (string, error) {
	if runtime.GOOS == "windows" {
		return AutoStartWindows(appPath)
	}
	if runtime.GOOS == "darwin" {
		return AutoStartDarwin(name, appPath)
	}
	//Darwin ?
	return "", nil
}

func AutoStartWindows(appPath string) (string, error) {
	userProfile := "USERPROFILE"
	userProfileFolder := os.Getenv(userProfile)
	if userProfile == "" {
		return "", fmt.Errorf("%s: %w", userProfile, ErrNoUserProfile)
	}
	appName := filepath.Base(appPath)
	fileName := strings.TrimSuffix(appName, filepath.Ext(appName))
	startupLinkPath := filepath.Join(userProfileFolder, "AppData", "Roaming", "Microsoft", "Windows", "Start Menu", "Programs", "Startup", fileName+".lnk")
	if err := makeLink(appPath, startupLinkPath, false); err != nil {
		return "", err
	}
	return startupLinkPath, nil
}

var plistTemplate = `
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
          <key>RunAtLoad</key>
          <true/>
          <key>Label</key>
          <string>%s</string>
          <key>Program</key>
          <string>%s</string>
          <key>KeepAlive</key>
          <true/>
</dict>
</plist>
`

func AutoStartDarwin(name string, path string) (string, error) {
	userProfile := "HOME"
	userProfileFolder := os.Getenv(userProfile)
	if userProfileFolder == "" {
		return "", fmt.Errorf("%s: %w", userProfile, ErrNoUserProfile)
	}
	folder := "Library/LaunchAgents"
	fullPath := filepath.Join(userProfileFolder, folder, name+".plist")
	f, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := fmt.Fprintf(f, plistTemplate, name, path); err != nil {
		return "", err
	}
	return fullPath, nil
}
