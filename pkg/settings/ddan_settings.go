/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

ddan_settings.go

Analyzer settings widgets
*/
package settings

import (
	"context"
	"errors"
	"image/color"
	"net/url"
	"sandboxer/pkg/config"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/mpkondrashin/ddan"
)

type DDAn struct {
	conf *config.DDAn

	ddanURLEntry       *widget.Entry
	ddanAPIKeyEntry    *widget.Entry
	ddanIgnoreTLSCheck *widget.Check
	ddanTest           *canvas.Text //  *widget.Label
	cancelTestDDAn     context.CancelFunc
}

func NewDDAnSettings(conf *config.DDAn) *DDAn {
	return &DDAn{
		conf: conf,
	}
}

func (s *DDAn) Widget() fyne.CanvasObject {

	s.ddanURLEntry = widget.NewEntry()
	s.ddanURLEntry.SetText(s.conf.GetURL())
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

	s.ddanTest = canvas.NewText("", color.Black)
	//s.ddanTest = widget.NewLabel("")
	//s.ddanTest.Truncation = fyne.TextTruncateEllipsis

	//stateText := canvas.NewText(tsk.GetChannel(), tsk.RiskLevel.Color())
	//stateText.TextStyle = fyne.TextStyle{Bold: tsk.Active}
	ddanForm := widget.NewForm(urlFormItem, apiKeyFormItem, ignoreTLSFormItem)
	return container.NewVBox(ddanForm, container.NewHScroll(s.ddanTest))
}

func (s *DDAn) Update() {
	s.TestAnalyzer()
}

/*
const MaxLength = 64

	func LimitLength(s string) string {
		logging.Errorf("DDAn Connection: %s", s)
		if len(s) < MaxLength {
			return "Error: " + s
		}
		return "Error: ..." + s[len(s)-MaxLength+7:]
	}
*/

func (s *DDAn) SetMessageError(message string) {
	s.ddanTest.Text = "Error: " + message
	s.ddanTest.Color = color.RGBA{255, 0, 0, 255}
	s.ddanTest.Refresh()
}

func (s *DDAn) SetMessageOk(message string) {
	s.ddanTest.Text = message
	s.ddanTest.Color = color.Black
	s.ddanTest.Refresh()
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
		s.SetMessageOk("Checking connection...")

		u, err := url.Parse(s.GetDDAnURL())
		if err != nil {
			s.SetMessageError(err.Error())
			return
		}
		apiKey := strings.TrimSpace(s.ddanAPIKeyEntry.Text)
		analyzer := ddan.NewClient(s.conf.GetProductName(), s.conf.GetHostname()).
			SetAnalyzer(u, apiKey, s.ddanIgnoreTLSCheck.Checked)
		analyzer.SetProtocolVersion(s.conf.GetProtocolVersion())

		ctxTimeout, cancelTimeout := context.WithTimeout(ctx, 5*time.Second)
		defer cancelTimeout()
		err = analyzer.TestConnection(ctxTimeout)
		if err != nil {
			if !errors.Is(err, context.Canceled) {
				if errors.Is(err, context.DeadlineExceeded) {
					s.SetMessageError("Connection timed out")
					//s.ddanTest.SetText("Connection timed out")
				} else {
					s.SetMessageError(err.Error())
				}
			}
		} else {
			s.SetMessageOk("Connection is Ok")
			//s.ddanTest.SetText("Connection is Ok")
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
	s.conf.SetURL(s.GetDDAnURL())
	apiKey := strings.TrimSpace(s.ddanAPIKeyEntry.Text)
	if apiKey == "" {
		return errors.New("Analyzer API Key is empty")
	}
	s.conf.APIKey = apiKey
	s.conf.IgnoreTLSErrors = s.ddanIgnoreTLSCheck.Checked
	return nil
}
