package sandbox

import "errors"

var (
	ErrUnsupported = errors.New("unsupported")
)

var (
	//ErrNotReady = errors.New("not ready")
	ErrNotFound = errors.New("not found")
	ErrError    = errors.New("error")
)

type Sandbox interface {
	SubmitURL(url string) (string, error)
	SubmitFile(filePath string) (string, error)
	//HaveFinished(id string) (bool, error)
	GetResult(id string) (RiskLevel, string, error)
	GetReport(id string, filePath string) error
	GetInvestigation(id string, filePath string) error
}

/*
V1 Submit(filePath string) (string, error) {

}
Submit - GetStatus -      Get Result

DDan
Submit - GetBriefReport - [GetBriefReport - GetReport]

V1
Submit - [GetStatus - Get Result]

DDan
Submit - [GetBriefReport - GetReport]


*/
