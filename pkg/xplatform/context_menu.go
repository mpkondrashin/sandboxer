package xplatform

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func ExtendContextMenu(appName, appPath string) (string, error) {
	if runtime.GOOS == "windows" {
		return ExtendContextMenuWindows(appName, appPath)
	}
	//Darwin ?
	return "", nil
}

func ExtendContextMenuWindows(appName, appPath string) (string, error) {
	appData := "APPDATA"
	userProfile := os.Getenv(appData)
	if userProfile == "" {
		return "", fmt.Errorf("%s: %w", appData, ErrNoUserProfile)
	}
	linkName := appName + ".lnk"
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
