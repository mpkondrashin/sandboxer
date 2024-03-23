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

var (
	ErrUnknownAuthType   = errors.New("unknown auth type")
	ErrMissingScheme     = errors.New("missing proxy scheme")
	ErrMissingURL        = errors.New("missing proxy username")
	ErrMissingUsername   = errors.New("missing proxy username")
	ErrMissingPassword   = errors.New("missing proxy password")
	ErrMissingNTLMDomain = errors.New("missing NTLM domain")
)

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
	if a < 0 || a >= AuthTypeNTLM {
		return nil, ErrUnknownAuthType
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

type YAMLURL struct {
	*url.URL
}

// MarshalYAML implements the Marshaler interface of the yaml.v3 package for YAMLURL.
func (u YAMLURL) MarshalYAML() (interface{}, error) {
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

type Proxy struct {
	mx        sync.RWMutex `gsetter:"-"`
	Active    bool
	Type      AuthType
	URL       YAMLURL
	Username  string
	Password  string
	Domain    string
	Timeout   time.Duration
	KeepAlive time.Duration
}

func NewProxy(URL *url.URL) *Proxy {
	return &Proxy{
		Active:    false,
		Type:      AuthTypeNone,
		URL:       YAMLURL{URL},
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
	p.URL = newProxy.URL
	p.Type = newProxy.Type
	p.Username = newProxy.Username
	p.Password = newProxy.Password
	p.Domain = newProxy.Domain
}

func (p *Proxy) BasicAuth(Username string, Password string) *Proxy {
	p.Type = AuthTypeBasic
	p.Username = Username
	p.Password = Password
	return p
}

func (p *Proxy) NTLMAuth(Username string, Password string, Domain string) *Proxy {
	p.Type = AuthTypeNTLM
	p.Username = Username
	p.Password = Password
	p.Domain = Domain
	return p
}

func (p *Proxy) Modifier() (func(*http.Transport), error) {
	if !p.Active {
		return NullTransportModifier, nil
	}
	if p.URL.URL.Scheme == "" {
		return nil, ErrMissingScheme
	}
	if p.URL.URL.Host == "" {
		return nil, ErrMissingURL
	}
	if p.Type == AuthTypeNone {
		return p.TransportNoAuth, nil
	}
	if p.Username == "" {
		return nil, ErrMissingUsername
	}
	if p.Password == "" {
		return nil, ErrMissingPassword
	}
	if p.Type == AuthTypeBasic {
		return p.TransportBasic, nil
	}
	if p.Domain == "" {
		return nil, ErrMissingNTLMDomain
	}
	if p.Type == AuthTypeNTLM {
		return p.TransportNTLM, nil
	}
	return nil, ErrUnknownAuthType
}

/*
func (p *Proxy) ChangeTransport(t *http.Transport) {

	switch p.Type {
	default:
		fallthrough
	case AuthTypeNone:
		p.TransportNoAuth(t)
	case AuthTypeBasic:
		p.TransportBasic(t)
	case AuthTypeNTLM:
		p.TransportNTLM(t)
	}
}*/

func (p *Proxy) TransportNoAuth(t *http.Transport) {
	t.Proxy = http.ProxyURL(p.URL.URL)
}

func (p *Proxy) TransportNTLM(t *http.Transport) {
	dialer := &net.Dialer{
		Timeout:   p.Timeout,
		KeepAlive: p.KeepAlive,
	}
	ntlmDialContext := ntlm.NewNTLMProxyDialContext(dialer, *p.URL.URL, p.Username, p.Password, p.Domain, nil)
	t.Proxy = nil
	t.DialContext = ntlmDialContext

}

func (p *Proxy) TransportBasic(t *http.Transport) {
	u := *p.URL.URL
	u.User = url.UserPassword(p.Username, p.Password)
	t.Proxy = http.ProxyURL(&u)
}

func NullTransportModifier(*http.Transport) {
}

/*


func (p *Proxy) ConfigureProxy() (TransportModifier, error) {
	if p.URL.URL == nil {
		return DummyTransportModifier, nil
	}
	proxy := config.NewProxy(p.URL.URL)
	if p.Username == "" {
		return proxy,
	} viper.GetString(flagProxyUser) != "" {
			if viper.GetString(flagProxyPassword) == "" {
				log.Fatal("missing proxy password")
			}
			if viper.GetString(flagProxyDomain) != "" {
				log.Println("Use NTLM proxy auth")
				proxy.NTLMAuth(
					viper.GetString(flagProxyUser),
					viper.GetString(flagProxyPassword),
					viper.GetString(flagProxyDomain),
				)
			} else {
				log.Println("Use basic proxy auth")
				proxy.BasicAuth(
					viper.GetString(flagProxyUser),
					viper.GetString(flagProxyPassword),
				)
			}
		}
		c.visionOne.AddTransportModifier(proxy.GetModifier())
	}

}
*/
