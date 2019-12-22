package service

import (
	"errors"

	"github.com/asaskevich/govalidator"
	"github.com/herwonowr/slackhell/internal/repository"
)

func (s *internalService) GetSessionByID(id string) (*repository.Session, error) {
	if govalidator.IsNull(id) {
		return nil, errors.New("invalid session id")
	}

	return s.repository.GetSessionByID(id)
}

func (s *internalService) CreateSession(id string, key string) error {
	if govalidator.IsNull(id) {
		return errors.New("invalid session id")
	}

	if govalidator.IsNull(key) {
		return errors.New("invalid session key")
	}

	session, err := s.repository.GetSessionByID(id)
	if err != nil && err.Error() != "key not found" {
		return err
	}

	if session != nil {
		return errors.New("session already exists")
	}

	sessionData := &repository.Session{
		ID:  id,
		Key: key,
	}

	err = s.repository.CreateSession(sessionData)
	if err != nil {
		return err
	}

	return nil
}

func (s *internalService) PutSession(id string, key string) error {
	if govalidator.IsNull(id) {
		return errors.New("invalid session id")
	}

	if govalidator.IsNull(key) {
		return errors.New("invalid session key")
	}

	user, err := s.repository.GetAccountBySlackID(id)
	if err != nil {
		return err
	}

	shellcode, err := s.repository.GetShellcodeByKey(key)
	if err != nil {
		return err
	}

	session := &repository.Session{
		ID:  user.SlackID,
		Key: shellcode.ShellKey,
	}
	if err := s.repository.PutSession(session); err != nil {
		return err
	}

	return nil
}

func (s *internalService) DeleteSession(id string) error {
	if govalidator.IsNull(id) {
		return errors.New("invalid session id")
	}

	return s.repository.DeleteSession(id)
}
