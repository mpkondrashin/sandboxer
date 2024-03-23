package config

import "time"

func (s *Proxy) GetActive() bool {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.Active
}

func (s *Proxy) SetActive(value bool ) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.Active = value
}

func (s *Proxy) GetType() AuthType {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.Type
}

func (s *Proxy) SetType(value AuthType ) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.Type = value
}

func (s *Proxy) GetURL() YAMLURL {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.URL
}

func (s *Proxy) SetURL(value YAMLURL ) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.URL = value
}

func (s *Proxy) GetUsername() string {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.Username
}

func (s *Proxy) SetUsername(value string ) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.Username = value
}

func (s *Proxy) GetPassword() string {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.Password
}

func (s *Proxy) SetPassword(value string ) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.Password = value
}

func (s *Proxy) GetDomain() string {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.Domain
}

func (s *Proxy) SetDomain(value string ) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.Domain = value
}

func (s *Proxy) GetTimeout() time.Duration {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.Timeout
}

func (s *Proxy) SetTimeout(value time.Duration ) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.Timeout = value
}

func (s *Proxy) GetKeepAlive() time.Duration {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.KeepAlive
}

func (s *Proxy) SetKeepAlive(value time.Duration ) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.KeepAlive = value
}

