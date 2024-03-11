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
}

func NewVisionOne(domain, token string) *VisionOne {
	return &VisionOne{
		Domain: domain,
		Token:  token,
	}
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
	return vone.NewVOne(domain, token), nil
}
