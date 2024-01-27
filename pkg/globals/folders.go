package globals

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sandboxer/pkg/logging"
	"strings"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

var (
	ErrNoUserProfile = errors.New("missing environment variable")
	ErrUnsupportedOS = errors.New("unsupported OS")
)

func UserDataFolder() (string, error) {
	if runtime.GOOS == "windows" {
		return userDataFolder("APPDATA", AppID, "")
	}
	if runtime.GOOS == "darwin" {
		return userDataFolder("HOME", "Library/Application Support", AppID)
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

func InstallFolder() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("PROGRAMFILES")
	}
	if runtime.GOOS == "darwin" {
		return "/Applications"
	}
	return ""
}

func ConfigurationFilePath() (string, error) {
	folder, err := UserDataFolder()
	if err != nil {
		return "", err
	}
	return filepath.Join(folder, ConfigFileName), nil
}

func LogsFolder() (string, error) {
	folder, err := UserDataFolder()
	if err != nil {
		return "", err
	}
	return filepath.Join(folder, "logs"), nil
}

func ExtendContextMenu(appPath string) (string, error) {
	if runtime.GOOS == "windows" {
		return ExtendContextMenuWindows(appPath)
	}
	//Darwin ?
	return "", nil
}

func ExtendContextMenuWindows(appPath string) (string, error) {
	appData := "APPDATA"
	userProfile := os.Getenv(appData)
	if userProfile == "" {
		return "", fmt.Errorf("%s: %w", appData, ErrNoUserProfile)
	}
	linkName := AppName + ".lnk"
	linkPath := filepath.Join(userProfile, "Microsoft", "Windows", "SendTo", linkName)
	_ = os.Remove(linkPath)
	if err := makeLink(appPath, linkPath); err != nil {
		return "", err
	}
	return linkPath, nil
}

func makeLink(src, dst string) error {
	ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY)
	oleShellObject, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return err
	}
	defer oleShellObject.Release()
	wshell, err := oleShellObject.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return err
	}
	defer wshell.Release()
	cs, err := oleutil.CallMethod(wshell, "CreateShortcut", dst)
	if err != nil {
		return err
	}
	idispatch := cs.ToIDispatch()
	if _, err := oleutil.PutProperty(idispatch, "TargetPath", src); err != nil {
		return err
	}
	if _, err := oleutil.CallMethod(idispatch, "Save"); err != nil {
		return err
	}
	return nil
}

func AutoStart(appPath string) (string, error) {
	if runtime.GOOS == "windows" {
		return AutoStartWindows(appPath)
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
	if err := makeLink(appPath, startupLinkPath); err != nil {
		return "", err
	}
	return startupLinkPath, nil
}

func PidFilePath() (string, error) {
	folder, err := UserDataFolder()
	if err != nil {
		return "", err
	}
	return filepath.Join(folder, Name+".pid"), nil
}

func SetupLogging(logFileName string) (func(), error) {
	logging.SetLevel(logging.DEBUG)
	//      logFileName := fmt.Sprintf("setup_%s.log", time.Now().Format("20060102_150405"))
	logFolder, err := LogsFolder()
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(logFolder, 0700); err != nil {
		return nil, err
	}
	close := logging.NewFileLog(logFolder, logFileName)
	return func() {
		logging.Infof("Close Logging")
		close()
	}, nil
}
