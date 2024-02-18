/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

folders.go

Various folders that used through whole project
*/
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
	"github.com/virtuald/go-paniclog"
)

const (
	tasksFolder = "tasks"
	logsFolder  = "logs"
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
	return filepath.Join(folder, logsFolder), nil
}

func TasksFolder() (string, error) {
	folder, err := UserDataFolder()
	if err != nil {
		return "", err
	}
	return filepath.Join(folder, tasksFolder), nil
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

func AutoStartDarwin(appPath string) (string, error) {

	appleScript := `set theScript to "/Applications/Sandboxer.app/Contents/MacOS/submit $1" -- the shell script

tell application "Automator"
   set actionID to Automator action id "com.apple.RunShellScript"
   tell (make new workflow)
      add actionID to it -- add to the end of the workflow
      tell last Automator action
         set value of setting "inputMethod" to 1 -- arguments menu
         set value of setting "COMMAND_STRING" to theScript
      end tell
   end tell
   activate
   # tell application "System Events" to keystroke return -- default workflow
end tell

set theScript to "/Applications/Sandboxer.app/Contents/MacOS/submit $1" -- the shell script

tell application "Automator"
	set actionID to Automator action id "com.apple.RunShellScript"
	tell (make new workflow)
		add actionID to it -- add to the end of the workflow
		tell last Automator action
			set value of setting "inputMethod" to 1 -- arguments menu
			set value of setting "COMMAND_STRING" to theScript
		end tell
		set receives to make new receive workflow item at beginning of workflows
		tell receives
			set receives to {alias}
		end tell
	end tell
	activate
	# tell application "System Events" to keystroke return -- default workflow
end tell

`
	_ = appleScript
	return "", nil
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
	logFolder, err := LogsFolder()
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(logFolder, 0755); err != nil {
		return nil, err
	}
	logging.SetLevel(logging.DEBUG)
	file, err := logging.OpenRotated(logFolder, logFileName, 0644, MaxLogFileSize, LogsKeep)
	if err != nil {
		return nil, err
	}
	paniclog.RedirectStderr(file.File)
	logging.SetLogger(logging.NewFileLogger(file))
	return func() {
		logging.Infof("Close Logging")
		file.Close()
	}, nil
}

func DownloadsFolder() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("USERPROFILE"), "Downloads")
	}
	return filepath.Join(os.Getenv("HOME"), "Downloads")
}
