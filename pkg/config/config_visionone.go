package config

import (
	"errors"
	"sync"

	"github.com/mpkondrashin/vone"
)

type VisionOne struct {
	mx     sync.RWMutex `gsetter:"-"`
	Token  string       `yaml:"token"`
	Domain string       `yaml:"domain"`
	Proxy  *Proxy       `yaml:"-" gsetter:"-"`
}

func NewVisionOne(domain, token string) *VisionOne {
	return &VisionOne{
		Domain: domain,
		Token:  token,
	}
}

func (v *VisionOne) Update(newVOne *VisionOne) {
	v.mx.Lock()
	defer v.mx.Unlock()
	newVOne.mx.RLock()
	defer newVOne.mx.RUnlock()
	v.Token = newVOne.Token
	v.Domain = newVOne.Domain
}

func (s *VisionOne) VisionOneSandbox() (*vone.VOne, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	token := s.Token
	if token == "" {
		return nil, errors.New("token is not set")
	}
	domain := s.Domain
	if domain == "" {
		return nil, errors.New("domain is not set")
	}
	v := vone.NewVOne(domain, token)
	if s.Proxy == nil {
		return v, nil
	}
	modifier, err := s.Proxy.Modifier()
	if err != nil {
		return nil, err
	}
	v.AddTransportModifier(modifier)
	return v, nil
}
