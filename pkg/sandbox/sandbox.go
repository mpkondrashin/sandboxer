/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

sandbox.go

Unified sandbox interface
*/
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
	GetResult(id string) (RiskLevel, string, error)
	GetReport(id string, filePath string) error
	GetInvestigation(id string, filePath string) error
}
