package settings

import (
	"context"
	"errors"
	"sandboxer/pkg/config"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/mpkondrashin/vone"
)

const ErrorDomain = "Select"

type VisionOne struct {
	conf *config.VisionOne

	tokenEntry *widget.Entry
	//domainLabel      *widget.Label
	//chooseLabel      *widget.Label
	visionOneDomains *widget.Select
	cancelDetect     context.CancelFunc
}

func NewVisionOne(conf *config.VisionOne) *VisionOne {
	return &VisionOne{
		conf: conf,
	}
}

func (s *VisionOne) Widget() fyne.CanvasObject {

	s.tokenEntry = widget.NewMultiLineEntry()
	s.tokenEntry.SetText(s.conf.Token)
	s.tokenEntry.Wrapping = fyne.TextWrapBreak
	s.tokenEntry.OnChanged = s.DetectDomain
	tokenFormItem := widget.NewFormItem("Token:", s.tokenEntry)
	tokenFormItem.HintText = "Go to Administrator -> API Keys"
	//roleHint := "Go to Administration -> User Roles -> Permissions -> Threat Intelligence -> Sandbox Analysis -> \"View, filter, and search\" and\"Submit object\""
	// apiKeyHitn := "Go to Administration -> API Keys -> Add API Key"

	//s.domainLabel = widget.NewLabel(s.conf.Domain)
	//domainFormItem := widget.NewFormItem("Domain:", s.domainLabel)
	//s.chooseLabel = widget.NewLabel("Choose Vision One Domain:")
	domains := []string{ErrorDomain}
	for _, d := range vone.RegionalDomains {
		//l := fmt.Sprintf("%s (%s)", d.Region, d.Domain)
		domains = append(domains, d.Domain)
	}

	s.visionOneDomains = widget.NewSelect(domains, nil)
	if s.conf.Domain != "" {
		s.visionOneDomains.SetSelected(s.conf.Domain)
	}
	domainFormItem := widget.NewFormItem("Domain:", s.visionOneDomains)
	optionsForm := widget.NewForm(
		tokenFormItem,
		domainFormItem,
	)
	return optionsForm
}

func (s *VisionOne) Update() {
	s.DetectDomain(s.tokenEntry.Text)
}

func (s *VisionOne) DetectDomain(token string) {
	go func() {
		if s.cancelDetect != nil {
			s.cancelDetect()
		}
		var ctx context.Context
		ctx, s.cancelDetect = context.WithCancel(context.TODO())
		defer func() {
			if s.cancelDetect != nil {
				s.cancelDetect()
			}
			s.cancelDetect = nil
		}()
		domain := vone.DetectVisionOneDomain(ctx, token)
		if domain != "" {
			s.visionOneDomains.SetSelected(domain)
		}
	}()
}

func (s *VisionOne) Aquire() error {
	s.conf.Token = strings.TrimSpace(s.tokenEntry.Text)
	if s.conf.Token == "" {
		return errors.New("Vision One Token is empty")
	}
	if s.visionOneDomains.Selected == "" {
		return errors.New("Vision One Domain is not selected")
	}
	s.conf.Domain = s.visionOneDomains.Selected
	return nil
}
