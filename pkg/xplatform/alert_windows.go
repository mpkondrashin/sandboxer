//go:build windows

package xplatform

import "github.com/go-toast/toast"

func Alert(title, subtitle, message, iconPath string) error {
	//globals.folder
	notification := toast.Notification{
		AppID:   title,
		Title:   subtitle,
		Message: message,
		Icon:    iconPath, // This file must exist (remove this line if it doesn't)
		/*Actions: []toast.Action{
			{"protocol", "I'm a button", ""},
			{"protocol", "Me too!", ""},
		},*/
	}
	return notification.Push()
}
