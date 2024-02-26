package xplatform

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

var (
	ErrUnsupportedOS = errors.New("unsupported OS")
	ErrNoUserProfile = errors.New("missing environment variable")
)

func InstallFolder() string {
	switch runtime.GOOS {
	case "windows":
		return os.Getenv("PROGRAMFILES")
	case "darwin":
		return "/Applications"
	default:
		return ""
	}
}

func UserDataFolder(appID string) (string, error) {
	if runtime.GOOS == "windows" {
		return userDataFolder("APPDATA", appID, "")
	}
	if runtime.GOOS == "darwin" {
		return userDataFolder("HOME", "Library/Application Support", appID)
	}
	return "", fmt.Errorf("%s: %w", runtime.GOOS, ErrUnsupportedOS)
}

func userDataFolder(profileVariable string, folder string, subfolder string) (string, error) {
	userProfile := os.Getenv(profileVariable)
	if userProfile == "" {
		return "", fmt.Errorf("%s: %w", profileVariable, ErrNoUserProfile)
	}
	return filepath.Join(userProfile, folder, subfolder), nil
}

func DownloadsFolder() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("USERPROFILE"), "Downloads")
	}
	return filepath.Join(os.Getenv("HOME"), "Downloads")
}

func ExecutablePath(folder string, appName string, name string) (string, error) {
	if runtime.GOOS == "windows" {
		return filepath.Join(folder, appName, name+".exe"), nil
	}
	if runtime.GOOS == "darwin" {
		return fmt.Sprintf("%s/%s.app/Contents/MacOS/%s", folder, appName, name), nil
	}
	return "", fmt.Errorf("%s: %W", runtime.GOOS, ErrUnsupportedOS)
}

func ExecutableName(name string) string {
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	return name
}
