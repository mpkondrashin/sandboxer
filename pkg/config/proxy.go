package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/launchdarkly/go-ntlm-proxy-auth"
	"gopkg.in/yaml.v3"
)

func AddTransportModifier(pModifier *func(*http.Transport), modifier func(*http.Transport)) {
	if *pModifier == nil {
		*pModifier = modifier
		return
	}
	original := *pModifier
	*pModifier = func(t *http.Transport) {
		original(t)
		modifier(t)
	}
}

var (
	ErrUnknownAuthType  = errors.New("unknown auth type")
	ErrUnknownProxyType = errors.New("unknown proxy type")
	//ErrMissingScheme     = errors.New("missing proxy scheme")
	ErrMissingAddress    = errors.New("missing address")
	ErrMissingPort       = errors.New("missing port")
	ErrMissingUsername   = errors.New("missing proxy username")
	ErrMissingPassword   = errors.New("missing proxy password")
	ErrMissingNTLMDomain = errors.New("missing NTLM domain")
)

type AuthType int

const (
	AuthTypeNone AuthType = iota
	AuthTypeBasic
	AuthTypeNTLM
)

var AuthTypeString = []string{
	"None",
	"Basic",
	"NTLM",
}

func (r AuthType) String() string {
	return AuthTypeString[r]
}

func AuthTypeFromString(s string) (AuthType, error) {
	for i, t := range AuthTypeString {
		if strings.EqualFold(t, s) {
			return AuthType(i), nil
		}
	}
	return 0, fmt.Errorf("%w: %s", ErrUnknownAuthType, s)
}

// UnmarshalJSON implements the Unmarshaler interface of the json package for AuthType.
func (a *AuthType) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	authType, err := AuthTypeFromString(v)
	if err != nil {
		return err
	}
	*a = authType
	return nil
}

// MarshalJSON implements the Marshaler interface of the json package for AuthType.
func (a AuthType) MarshalJSON() ([]byte, error) {
	if a < 0 || a > AuthTypeNTLM {
		return nil, fmt.Errorf("%d: %w", a, ErrUnknownAuthType)
	}
	return []byte(fmt.Sprintf("\"%s\"", a.String())), nil
}

// MarshalYAML implements the Marshaler interface of the yaml.v3 package for AuthType.
func (s AuthType) MarshalYAML() (interface{}, error) {
	return s.String(), nil
}

// UnmarshalYAML implements the Unmarshaler interface of the yaml.v3 package for AuthType.
func (a *AuthType) UnmarshalYAML(value *yaml.Node) error {
	var v string
	err := value.Decode(&v)
	if err != nil {
		return err
	}
	authType, err := AuthTypeFromString(v)
	if err != nil {
		return err
	}
	*a = authType
	return nil
}

// YAMLURL
/*
type YAMLURL struct {
	*url.URL
}

// MarshalYAML implements the Marshaler interface of the yaml.v3 package for YAMLURL.
func (u YAMLURL) MarshalYAML() (interface{}, error) {
	if u.URL == nil {
		return "", nil
	}
	return u.URL.String(), nil
}

func (y *YAMLURL) UnmarshalYAML(value *yaml.Node) error {
	var v string
	err := value.Decode(&v)
	if err != nil {
		return err
	}
	url, err := url.Parse(v)
	if err != nil {
		return err
	}
	y.URL = url
	return nil
}
*/
type Proxy struct {
	mx        sync.RWMutex `gsetter:"-"`
	Active    bool
	Address   string
	Port      int
	AuthType  AuthType
	Username  string
	Password  string
	Domain    string
	Timeout   time.Duration
	KeepAlive time.Duration
}

func NewProxy() *Proxy {
	return &Proxy{
		Active:    false,
		AuthType:  AuthTypeNone,
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
}

func (p *Proxy) Update(newProxy *Proxy) {
	p.mx.Lock()
	defer p.mx.Unlock()
	newProxy.mx.RLock()
	defer newProxy.mx.RUnlock()
	p.Active = newProxy.Active
	p.Address = newProxy.Address
	p.Port = newProxy.Port
	p.AuthType = newProxy.AuthType
	p.Username = newProxy.Username
	p.Password = newProxy.Password
	p.Domain = newProxy.Domain
}

func (p *Proxy) Modifier() (func(*http.Transport), error) {
	if !p.Active {
		return NullTransportModifier, nil
	}
	if p.Address == "" {
		return nil, ErrMissingAddress
	}
	if p.Port == 0 {
		return nil, ErrMissingPort
	}
	if p.AuthType == AuthTypeNone {
		return p.TransportNoAuth, nil
	}
	if p.Username == "" {
		return nil, ErrMissingUsername
	}
	if p.Password == "" {
		return nil, ErrMissingPassword
	}
	if p.AuthType == AuthTypeBasic {
		return p.TransportBasic, nil
	}
	if p.Domain == "" {
		return nil, ErrMissingNTLMDomain
	}
	if p.AuthType == AuthTypeNTLM {
		return p.TransportNTLM, nil
	}
	return nil, ErrUnknownAuthType
}

func NullTransportModifier(*http.Transport) {
}

func (p *Proxy) TransportNoAuth(t *http.Transport) {
	u := &url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", p.Address, p.Port),
	}
	t.Proxy = http.ProxyURL(u)
}

func (p *Proxy) TransportBasic(t *http.Transport) {
	u := &url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", p.Address, p.Port),
		User:   url.UserPassword(p.Username, p.Password),
	}
	t.Proxy = http.ProxyURL(u)
}

func (p *Proxy) TransportNTLM(t *http.Transport) {
	dialer := &net.Dialer{
		Timeout:   p.Timeout,
		KeepAlive: p.KeepAlive,
	}
	u := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", p.Address, p.Port),
	}
	ntlmDialContext := ntlm.NewNTLMProxyDialContext(dialer, u, p.Username, p.Password, p.Domain, nil)
	t.Proxy = nil
	t.DialContext = ntlmDialContext

}
