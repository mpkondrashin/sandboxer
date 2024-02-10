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
	//SvcFileName    = "exam-ensvc.exe" // XXX only windows?
	//SvcName        = Name
	//SvcDisplayName = "Ex-amen Sandbox Submission Tool"
	//SvcDescription = "Submit files to Vision One sandbox"
)
