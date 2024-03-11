package config

import "time"

func (s *Configuration) GetfilePath() string {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.filePath
}

func (s *Configuration) SetfilePath(value string ) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.filePath = value
}

func (s *Configuration) GetVersion() string {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.Version
}

func (s *Configuration) SetVersion(value string ) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.Version = value
}

func (s *Configuration) GetSandboxType() SandboxType {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.SandboxType
}

func (s *Configuration) SetSandboxType(value SandboxType ) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.SandboxType = value
}

func (s *Configuration) GetFolder() string {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.Folder
}

func (s *Configuration) SetFolder(value string ) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.Folder = value
}

func (s *Configuration) GetIgnore() []string {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.Ignore
}

func (s *Configuration) SetIgnore(value []string ) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.Ignore = value
}

func (s *Configuration) GetSleep() time.Duration {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.Sleep
}

func (s *Configuration) SetSleep(value time.Duration ) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.Sleep = value
}

func (s *Configuration) GetPericulosum() string {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.Periculosum
}

func (s *Configuration) SetPericulosum(value string ) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.Periculosum = value
}

func (s *Configuration) GetShowPasswordHint() bool {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.ShowPasswordHint
}

func (s *Configuration) SetShowPasswordHint(value bool ) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.ShowPasswordHint = value
}

func (s *Configuration) GetTasksKeepDays() int {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.TasksKeepDays
}

func (s *Configuration) SetTasksKeepDays(value int ) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.TasksKeepDays = value
}

func (s *Configuration) GetShowNotifications() bool {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.ShowNotifications
}

func (s *Configuration) SetShowNotifications(value bool ) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.ShowNotifications = value
}

