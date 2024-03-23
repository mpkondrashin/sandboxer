/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

vone_settings.go

Vision One sandbox settings widgets
*/
package settings

import (
	"net/url"
	"sandboxer/pkg/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Proxy struct {
	Conf *config.Proxy

	activeCheck   *widget.Check
	urlEntry      *widget.Entry
	authTypeRadio *widget.RadioGroup
	usernameEntry *widget.Entry
	passwordEntry *widget.Entry
	domainEntry   *widget.Entry
	form          *widget.Form
	// cancelDetect     context.CancelFunc
}

func NewProxy(conf *config.Proxy) *Proxy {
	return &Proxy{
		Conf: conf,
	}
}

func (s *Proxy) Widget() fyne.CanvasObject {
	s.activeCheck = widget.NewCheck("Use proxy", func(check bool) {
		s.AuthTypeChange(s.authTypeRadio.Selected)
	})
	s.activeCheck.Checked = s.Conf.Active

	s.urlEntry = widget.NewEntry()
	if s.Conf.URL.URL != nil {
		s.urlEntry.SetText(s.Conf.URL.String())
	}
	urlFormItem := widget.NewFormItem("URL:", s.urlEntry)
	urlFormItem.HintText = "scheme://address:port, where scheme = socks, http, https"

	s.authTypeRadio = widget.NewRadioGroup(config.AuthTypeString, s.AuthTypeChange)
	s.authTypeRadio.Horizontal = true
	s.authTypeRadio.Required = true
	s.authTypeRadio.SetSelected(s.Conf.Type.String())
	authTypeFormItem := widget.NewFormItem("Auth Type:", s.authTypeRadio)

	s.usernameEntry = widget.NewEntry()
	s.usernameEntry.SetText(s.Conf.Username)
	usernameFormItem := widget.NewFormItem("Username:", s.usernameEntry)

	s.passwordEntry = widget.NewEntry()
	s.passwordEntry.SetText(s.Conf.Password)
	s.passwordEntry.Password = true
	passwordFormItem := widget.NewFormItem("Password:", s.passwordEntry)

	s.domainEntry = widget.NewEntry()
	s.domainEntry.SetText(s.Conf.Domain)
	domainFormItem := widget.NewFormItem("Domain:", s.domainEntry)

	s.form = widget.NewForm(
		urlFormItem,
		authTypeFormItem,
		usernameFormItem,
		passwordFormItem,
		domainFormItem,
	)
	s.AuthTypeChange(s.authTypeRadio.Selected)
	return container.NewVBox(
		s.activeCheck,
		s.form,
	)
}

func (s *Proxy) AuthTypeChange(choice string) {
	if s.activeCheck == nil {
		return
	}
	if s.form == nil {
		return
	}
	if s.usernameEntry == nil {
		return
	}
	if s.passwordEntry == nil {
		return
	}
	if s.domainEntry == nil {
		return
	}
	if !s.activeCheck.Checked {
		s.urlEntry.Disable()
		s.authTypeRadio.Disable()
		s.usernameEntry.Disable()
		s.passwordEntry.Disable()
		s.domainEntry.Disable()
		s.form.Refresh()
		return
	}
	s.authTypeRadio.Enable()
	s.urlEntry.Enable()
	switch choice {
	case config.AuthTypeNone.String():
		s.usernameEntry.Disable()
		s.passwordEntry.Disable()
		s.domainEntry.Disable()
	case config.AuthTypeBasic.String():
		s.usernameEntry.Enable()
		s.passwordEntry.Enable()
		s.domainEntry.Disable()
	case config.AuthTypeNTLM.String():
		s.usernameEntry.Enable()
		s.passwordEntry.Enable()
		s.domainEntry.Enable()
	}
	s.form.Refresh()
}

func (s *Proxy) Update() {
	//s.DetectDomain(s.tokenEntry.Text)
}

func (s *Proxy) Aquire() error {
	u, err := url.Parse(s.urlEntry.Text)
	if err != nil {
		return err
	}
	p := config.NewProxy(u)
	p.Active = s.activeCheck.Checked
	p.Type, err = config.AuthTypeFromString(s.authTypeRadio.Selected)
	if err != nil {
		return err
	}
	p.Username = s.usernameEntry.Text
	p.Password = s.passwordEntry.Text
	p.Domain = s.domainEntry.Text
	_, err = p.Modifier()
	if err != nil {
		return err
	}
	s.Conf.Update(p)
	return nil
}
