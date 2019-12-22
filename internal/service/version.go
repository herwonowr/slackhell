package service

func (s *internalService) GetVersion() (int, error) {
	version, err := s.repository.GetVersion()
	if err != nil {
		return 0, err
	}

	return version, nil
}
func (s *internalService) PutVersion(version int) error {
	return s.repository.PutVersion(version)
}
