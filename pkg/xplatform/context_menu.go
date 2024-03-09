/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

context_menu.go

Add extension to context menu for Finder/Explorer
*/

package xplatform

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func ExtendContextMenu(dryRun bool, appName, appPath string) (string, error) {
	if runtime.GOOS == "windows" {
		return ExtendContextMenuWindows(dryRun, appName, appPath)
	}
	//Darwin ?
	return "", nil
}

func ExtendContextMenuWindows(dryRun bool, appName, appPath string) (string, error) {
	appData := "APPDATA"
	userProfile := os.Getenv(appData)
	if userProfile == "" {
		return "", fmt.Errorf("%s: %w", appData, ErrNoUserProfile)
	}
	linkName := appName + ".lnk"
	linkPath := filepath.Join(userProfile, "Microsoft", "Windows", "SendTo", linkName)
	if dryRun {
		return linkPath, nil
	}
	_ = os.Remove(linkPath)
	if err := makeLink(appPath, linkPath, false); err != nil {
		return "", err
	}
	return linkPath, nil
}

func makeLink(src, dst string, asAdministrator bool) error {
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
	if asAdministrator {
		if err := runAsAdministrator(dst); err != nil {
			return nil
		}
	}
	return nil
}

// https://stackoverflow.com/questions/28997799/how-to-create-a-run-as-administrator-shortcut-using-powershell#:~:text=In%20short%2C%20you%20need%20to,This%20is%20the%20RunAsAdministrator%20flag.
func runAsAdministrator(path string) error {
	//In short, you need to read the .lnk file in as an array of bytes. Locate byte 21 (0x15) and change bit 6 (0x20) to 1. This is the RunAsAdministrator flag.
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	data[0x15] |= 0x20
	return os.WriteFile(path, data, 0644)
}
