/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

vone_settings.go

Vision One sandbox settings widgets
*/
package settings

import (
	"context"
	"sandboxer/pkg/config"
	"sandboxer/pkg/logging"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/mpkondrashin/vone"
)

const ErrorDomain = "Select"

type VisionOne struct {
	Conf *config.VisionOne

	tokenEntry *widget.Entry
	//domainLabel      *widget.Label
	//chooseLabel      *widget.Label
	visionOneDomains *widget.Select
	cancelDetect     context.CancelFunc
}

func NewVisionOne(conf *config.VisionOne) *VisionOne {
	return &VisionOne{
		Conf: conf,
	}
}

func (s *VisionOne) Widget() fyne.CanvasObject {

	s.tokenEntry = widget.NewMultiLineEntry()
	s.tokenEntry.SetText(s.Conf.GetToken())
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
	if s.Conf.GetDomain() != "" {
		s.visionOneDomains.SetSelected(s.Conf.GetDomain())
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
		logging.Debugf("Run Vision One domain detection")
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
	c := config.NewVisionOne(s.visionOneDomains.Selected, strings.TrimSpace(s.tokenEntry.Text))
	_, err := c.VisionOneSandbox()
	if err != nil {
		return err
	}
	s.Conf = c
	return nil
}
