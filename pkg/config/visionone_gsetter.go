package config

func (s *VisionOne) GetToken() string {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.Token
}

func (s *VisionOne) SetToken(value string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.Token = value
}

func (s *VisionOne) GetDomain() string {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.Domain
}

func (s *VisionOne) SetDomain(value string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.Domain = value
}
