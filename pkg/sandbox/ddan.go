/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

ddan.go

Sandbox for Deep Discovery Analyzer
*/
package sandbox

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mpkondrashin/ddan"
)

type DDAnSandbox struct {
	analyzer ddan.ClientInterface
}

var _ Sandbox = &DDAnSandbox{}

func NewDDAnSandbox(analyzer ddan.ClientInterface) *DDAnSandbox {
	return &DDAnSandbox{
		analyzer: analyzer,
	}
}

func (s *DDAnSandbox) SubmitURL(filePath string) (string, error) {
	return s.Submit(false, filePath)
}

func (s *DDAnSandbox) SubmitFile(filePath string) (string, error) {
	return s.Submit(true, filePath)
}

func CalculateStringHash(input string) string {
	hash := sha1.New()
	hash.Write([]byte(input))
	return hex.EncodeToString(hash.Sum(nil))
}

func (s *DDAnSandbox) Submit(file bool, content string) (string, error) {
	var sha1 string
	var err error
	if file {
		sha1, err = ddan.Hash(content)
		if err != nil {
			return "", err
		}
	} else {
		sha1 = CalculateStringHash(content)
	}
	sha1List, err := s.analyzer.CheckDuplicateSample(context.TODO(), []string{sha1}, 0)
	if err != nil {
		var apiErr *ddan.APIError
		if !errors.As(err, &apiErr) {
			return "", err
		}
		if apiErr.Response != ddan.ResponseNotRegistered {
			return "", err
		}
		err := s.analyzer.Register(context.TODO())
		if err != nil {
			return "", err
		}
		sha1List, err = s.analyzer.CheckDuplicateSample(context.TODO(), []string{sha1}, 0)
		if err != nil {
			return "", err
		}
	}
	if len(sha1List) == 1 {
		return sha1, nil
	}
	if file {
		err = s.analyzer.UploadSampleEx(context.TODO(), content, filepath.Base(content), sha1)
	} else {
		err = s.analyzer.UploadSampleURLs(context.TODO(), []string{content})
	}
	if err != nil {
		return "", err
	}
	return sha1, nil
}

func (s *DDAnSandbox) GetResult(id string) (RiskLevel, string, error) {
	briefReports, err := s.analyzer.GetBriefReport(context.TODO(), []string{id})
	if err != nil {
		return RiskLevelUnknown, "", fmt.Errorf("GetBriefReport: %w", err)
	}
	if len(briefReports.Reports) != 1 {
		return RiskLevelUnknown, "", fmt.Errorf("%s: %w: wrong brief report length", id, ErrError)
	}
	briefReport := briefReports.Reports[0]
	switch briefReport.SampleStatus {
	case ddan.StatusNotFound:
		return RiskLevelUnknown, "", fmt.Errorf("%s: %w", id, ErrNotFound)
	case ddan.StatusArrived:
		fallthrough
	case ddan.StatusProcessing:
		return RiskLevelNotReady, "", nil
	case ddan.StatusDone:

	case ddan.StatusError:
		return RiskLevelUnknown, "", fmt.Errorf("%s: %w", id, ErrError)
	case ddan.StatusTimeout:
		return RiskLevelError, "", fmt.Errorf("%s: %w: timeout", id, ErrError)
	default:
		return RiskLevelError, "", fmt.Errorf("%s: %w: %d: unknown status", id, ErrError, briefReport.SampleStatus)
	}
	switch briefReport.RiskLevel {
	case ddan.RatingUnsupported:
		return RiskLevelUnsupported, "", nil
	case ddan.RatingNoRiskFound:
		return RiskLevelNoRisk, "", nil
	}
	if briefReport.RiskLevel < 0 {
		return RiskLevelUnknown, "", fmt.Errorf("%s: %w: %v", id, ErrError, briefReport.RiskLevel)
	}
	reports, err := s.analyzer.GetReport(context.TODO(), id)
	if err != nil {
		return RiskLevelError, "", fmt.Errorf("GetReport(%s): %w", id, err)
	}
	if len(reports.FILEANALYZEREPORT) != 1 {
		return RiskLevelError, "", fmt.Errorf("%s: %w: wrong report length: %v", id, ErrError, reports)
	}
	VirusName := reports.FILEANALYZEREPORT[0].VirusName.Value

	switch briefReport.RiskLevel {
	case ddan.RatingLowRisk:
		return RiskLevelLow, VirusName, nil
	case ddan.RatingMediumRisk:
		return RiskLevelMedium, VirusName, nil
	case ddan.RatingHighRisk:
		return RiskLevelHigh, VirusName, nil
	default:
		return RiskLevelError, "", fmt.Errorf("GetBriefReport(%s): %d: %w", id, briefReport.RiskLevel, ErrUnknownRiskLevel)
	}
}

func (s *DDAnSandbox) GetReport(id string, filePath string) error {
	return GetFile(id, filePath, s.analyzer.GetPDFReport)
}

func (s *DDAnSandbox) GetInvestigation(id string, filePath string) error {
	return GetFile(id, filePath, s.analyzer.GetPackage)
}

func GetFile(id string, filePath string, apiCall func(context.Context, string) (io.Reader, error)) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	reader, err := apiCall(context.TODO(), id)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, reader)
	return err
}
