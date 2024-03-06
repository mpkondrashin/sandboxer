/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

globals.go

Global values
*/
package globals

const (
	AppName        = "Sandboxer"
	AppFolderName  = AppName
	Name           = "sandboxer"
	AppID          = "com.github.mpkondrashin." + Name
	ConfigFileName = Name + ".yaml"
	FIFOName       = Name + "_submit_fifo"
	MaxLogFileSize = 10_000_000
	LogsKeep       = 1
)
