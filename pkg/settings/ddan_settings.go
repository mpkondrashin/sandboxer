package settings

import (
	"context"
	"errors"
	"net/url"
	"sandboxer/pkg/config"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/mpkondrashin/ddan"
)

type DDAn struct {
	conf *config.DDAn

	ddanURLEntry       *widget.Entry
	ddanAPIKeyEntry    *widget.Entry
	ddanIgnoreTLSCheck *widget.Check
	ddanTest           *widget.Label
	cancelTestDDAn     context.CancelFunc
}

func NewDDAnSettings(conf *config.DDAn) *DDAn {
	return &DDAn{
		conf: conf,
	}
}

func (s *DDAn) Widget() fyne.CanvasObject {

	s.ddanURLEntry = widget.NewEntry()
	s.ddanURLEntry.SetText(s.conf.URL)
	s.ddanURLEntry.OnChanged = func(string) {
		s.TestAnalyzer()
	}
	urlFormItem := widget.NewFormItem("Address:", s.ddanURLEntry)
	urlFormItem.HintText = "DNS name or IP address"

	s.ddanAPIKeyEntry = widget.NewEntry()
	s.ddanAPIKeyEntry.SetText(s.conf.APIKey)
	s.ddanAPIKeyEntry.OnChanged = func(string) {
		s.TestAnalyzer()
	}
	apiKeyFormItem := widget.NewFormItem("API Key:", s.ddanAPIKeyEntry)
	apiKeyFormItem.HintText = "Go to Help -> About on Analyzer console"

	s.ddanIgnoreTLSCheck = widget.NewCheck("Ignore", nil)
	s.ddanIgnoreTLSCheck.SetChecked(s.conf.IgnoreTLSErrors)
	s.ddanIgnoreTLSCheck.OnChanged = func(bool) {
		s.TestAnalyzer()
	}
	ignoreTLSFormItem := widget.NewFormItem("TLS Errors: ", s.ddanIgnoreTLSCheck)

	s.ddanTest = widget.NewLabel("")

	ddanForm := widget.NewForm(urlFormItem, apiKeyFormItem, ignoreTLSFormItem)
	return container.NewVBox(ddanForm, s.ddanTest)
}

func (s *DDAn) Update() {
	s.TestAnalyzer()
}

func (s *DDAn) TestAnalyzer() {
	go func() {
		if s.cancelTestDDAn != nil {
			s.cancelTestDDAn()
		}
		var ctx context.Context
		ctx, s.cancelTestDDAn = context.WithCancel(context.TODO())
		defer func() {
			if s.cancelTestDDAn != nil {
				s.cancelTestDDAn()
			}
			s.cancelTestDDAn = nil
		}()
		s.ddanTest.SetText("Checking connection...")
		u, err := url.Parse(s.GetDDAnURL())
		if err != nil {
			s.ddanTest.SetText(err.Error())
			return
		}
		apiKey := strings.TrimSpace(s.ddanAPIKeyEntry.Text)
		analyzer := ddan.NewClient(s.conf.ProductName, s.conf.Hostname).
			SetAnalyzer(u, apiKey, s.ddanIgnoreTLSCheck.Checked)
		//if s.conf.ProtocolVersion != "" {
		//	log.Println("analyzer set version ", s.conf.ProtocolVersion)
		analyzer.SetProtocolVersion(s.conf.ProtocolVersion)
		//}
		ctxTimeout, cancelTimeout := context.WithTimeout(ctx, 5*time.Second)
		defer cancelTimeout()
		err = analyzer.TestConnection(ctxTimeout)
		if err != nil {
			if !errors.Is(err, context.Canceled) {
				if errors.Is(err, context.DeadlineExceeded) {
					s.ddanTest.SetText("Connection timed out")
				} else {
					s.ddanTest.SetText(err.Error())
				}
			}
		} else {
			s.ddanTest.SetText("Connection is Ok")
		}
	}()
}

func (s *DDAn) GetDDAnURL() (result string) {
	result = strings.TrimSpace(s.ddanURLEntry.Text)
	if strings.HasPrefix(result, "https://") {
		return
	}
	if strings.HasPrefix(result, "http://") {
		return
	}
	return "https://" + result
}

func (s *DDAn) Aquire() error {
	if strings.TrimSpace(s.ddanURLEntry.Text) == "" {
		return errors.New("Analyzer URL is empty")
	}
	s.conf.URL = s.GetDDAnURL()
	apiKey := strings.TrimSpace(s.ddanAPIKeyEntry.Text)
	if apiKey == "" {
		return errors.New("Analyzer API Key is empty")
	}
	s.conf.APIKey = apiKey
	s.conf.IgnoreTLSErrors = s.ddanIgnoreTLSCheck.Checked
	return nil
}
