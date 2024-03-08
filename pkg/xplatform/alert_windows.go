//go:build windows

package xplatform

import "github.com/go-toast/toast"

func Alert(appID, title, subtitle, message string) error {
	notification := toast.Notification{
		AppID:   appID,
		Title:   title,
		Message: subtitle + "\n" + message,
		//		Icon:    "go.png", // This file must exist (remove this line if it doesn't)
		/*Actions: []toast.Action{
			{"protocol", "I'm a button", ""},
			{"protocol", "Me too!", ""},
		},*/
	}
	return notification.Push()
}
