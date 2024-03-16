/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

vone.go

Sandbox for Vision One sandbox service
*/
package sandbox

import (
	"context"
	"fmt"
	"strings"

	"github.com/mpkondrashin/vone"
)

type VOneSandbox struct {
	vOne *vone.VOne
}

var _ Sandbox = &VOneSandbox{}

func NewVOneSandbox(vOne *vone.VOne) *VOneSandbox {
	return &VOneSandbox{
		vOne: vOne,
	}
}

func (s *VOneSandbox) SubmitURL(url string) (string, error) {
	f := s.vOne.SandboxSubmitURLs().AddURL(url)
	response, _, err := f.Do(context.TODO())
	if err != nil {
		return "", err
	}

	if len(response) != 1 {
		return "", fmt.Errorf("%s: %w: wrong response length: %v", url, ErrError, response)
	}
	return response[0].Body.ID, nil
}
func (s *VOneSandbox) SubmitFile(filePath string) (string, error) {
	f, err := s.vOne.SandboxSubmitFile().SetFilePath(filePath)
	if err != nil {
		return "", fmt.Errorf("Vision One ", err)
	}
	response, _, err := f.Do(context.TODO())
	if err != nil {
		return "", err
	}
	return response.ID, nil
}

func (s *VOneSandbox) GetResult(id string) (RiskLevel, string, error) {
	status, err := s.vOne.SandboxSubmissionStatus(id).Do(context.TODO())
	if err != nil {
		return RiskLevelUnknown, "", fmt.Errorf("SandboxSubmissionStatus(%s): %w", id, err)
	}
	switch status.Status {
	case vone.StatusSucceeded:
	case vone.StatusRunning:
		return RiskLevelNotReady, "", nil
	case vone.StatusFailed:
		if status.Error.Code == "Unsupported" {
			return RiskLevelUnsupported, "", fmt.Errorf("%w: %s", ErrUnsupported, status.Error.Message)
		}
		return RiskLevelError, "", fmt.Errorf("%s: %w: %s %s", id, ErrError, status.Error.Code, status.Error.Message)
	default:
		return RiskLevelError, "", fmt.Errorf("%v: %w", status, ErrUnknownRiskLevel)
	}
	results, err := s.vOne.SandboxAnalysisResults(id).Do(context.TODO())
	if err != nil {
		return RiskLevelUnknown, "", fmt.Errorf("SandboxAnalysisResults(%s): %w", id, err)
	}
	detectionName := strings.Join(results.DetectionNames, ", ")
	threatType := strings.Join(results.ThreatTypes, ", ")
	virusName := detectionName + threatType

	switch results.RiskLevel {
	case vone.RiskLevelNoRisk:
		return RiskLevelNoRisk, "", nil
	case vone.RiskLevelHigh:
		return RiskLevelHigh, virusName, nil
	case vone.RiskLevelMedium:
		return RiskLevelMedium, virusName, nil
	case vone.RiskLevelLow:
		return RiskLevelLow, virusName, nil
	default:
		return RiskLevelError, "", fmt.Errorf("%d: %w", results.RiskLevel, ErrUnknownRiskLevel)
	}
}

func (s *VOneSandbox) GetReport(id string, filePath string) error {
	return s.vOne.SandboxDownloadResults(id).Store(context.TODO(), filePath)
}
func (s *VOneSandbox) GetInvestigation(id string, filePath string) error {
	return s.vOne.SandboxInvestigationPackage(id).Store(context.TODO(), filePath)
}
