package config

func (s *DDAn) GetURL() string {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.URL
}

func (s *DDAn) SetURL(value string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.URL = value
}

func (s *DDAn) GetProtocolVersion() string {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.ProtocolVersion
}

func (s *DDAn) SetProtocolVersion(value string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.ProtocolVersion = value
}

func (s *DDAn) GetUserAgent() string {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.UserAgent
}

func (s *DDAn) SetUserAgent(value string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.UserAgent = value
}

func (s *DDAn) GetProductName() string {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.ProductName
}

func (s *DDAn) SetProductName(value string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.ProductName = value
}

func (s *DDAn) GetHostname() string {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.Hostname
}

func (s *DDAn) SetHostname(value string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.Hostname = value
}

func (s *DDAn) GetTempFolder() string {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.TempFolder
}

func (s *DDAn) SetTempFolder(value string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.TempFolder = value
}

func (s *DDAn) GetSourceID() string {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.SourceID
}

func (s *DDAn) SetSourceID(value string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.SourceID = value
}

func (s *DDAn) GetSourceName() string {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.SourceName
}

func (s *DDAn) SetSourceName(value string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.SourceName = value
}

func (s *DDAn) GetAPIKey() string {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.APIKey
}

func (s *DDAn) SetAPIKey(value string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.APIKey = value
}

func (s *DDAn) GetIgnoreTLSErrors() bool {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.IgnoreTLSErrors
}

func (s *DDAn) SetIgnoreTLSErrors(value bool) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.IgnoreTLSErrors = value
}

func (s *DDAn) GetClientUUID() string {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.ClientUUID
}

func (s *DDAn) SetClientUUID(value string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.ClientUUID = value
}
