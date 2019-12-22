package service

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/herwonowr/slackhell/internal/repository"
)

func (s *internalService) GetShellcodeByID(id int) (*repository.Shellcode, error) {
	if govalidator.IsNull(strconv.Itoa(id)) {
		return nil, errors.New("invalid shellcode id")
	}

	return s.repository.GetShellcodeByID(id)
}

func (s *internalService) GetShellcodeByKey(key string) (*repository.Shellcode, error) {
	if govalidator.IsNull(key) {
		return nil, errors.New("invalid shellcode key")
	}

	return s.repository.GetShellcodeByKey(key)
}

func (s *internalService) GetShellcodes() ([]repository.Shellcode, error) {
	return s.repository.GetShellcodes()
}

func (s *internalService) CreateShellcode(fileID string, shellkey string, ownerID string, ownerRealName string) error {
	if govalidator.IsNull(fileID) {
		return errors.New("invalid shellcode file id")
	}

	if govalidator.IsNull(shellkey) {
		return errors.New("invalid shellcode key")
	}

	if govalidator.IsNull(ownerID) {
		return errors.New("invalid shellcode owner")
	}

	aShellcodeKey, err := s.repository.GetShellcodeByKey(shellkey)
	if err != nil && err.Error() != "key not found" {
		return err
	}

	if aShellcodeKey != nil {
		return errors.New("shellcode with this key already exists")
	}

	shellOwner, err := s.repository.GetAccountBySlackID(ownerID)
	if err != nil && err.Error() == "account not found" {
		return errors.New("Owner not found")
	} else if err != nil {
		return err
	}

	shellcode := &repository.Shellcode{
		FileID:        fileID,
		ShellKey:      shellkey,
		OwnerID:       shellOwner.SlackID,
		OwnerRealName: shellOwner.SlackRealName,
	}

	err = s.repository.CreateShellcode(shellcode)
	if err != nil {
		return err
	}

	return nil
}

func (s *internalService) PutShellcode(key string, endpoint string, ownerID string) error {
	if govalidator.IsNull(key) {
		return errors.New("invalid shellcode key")
	}

	trimEndpoint := strings.TrimPrefix(endpoint, "<")
	cleanEndpoint := strings.TrimSuffix(trimEndpoint, ">")
	log.Println(cleanEndpoint)
	if !govalidator.IsURL(cleanEndpoint) {
		return errors.New("invalid shellcode endpoint")
	}

	if govalidator.IsNull(ownerID) {
		return errors.New("invalid shellcode owner id")
	}

	shellcode, err := s.repository.GetShellcodeByKey(key)
	if err != nil {
		return err
	}

	owner, err := s.repository.GetAccountBySlackID(ownerID)
	if err != nil {
		return err
	}

	if owner.Role != repository.AdminRole {
		if shellcode.OwnerID != ownerID {
			return errors.New("you don't have permission to update this shellcode")
		}
	}

	data := &repository.Shellcode{
		ID:            shellcode.ID,
		FileID:        shellcode.FileID,
		ShellKey:      shellcode.ShellKey,
		Endpoint:      cleanEndpoint,
		OwnerID:       shellcode.OwnerID,
		OwnerRealName: shellcode.OwnerRealName,
	}

	if err := s.repository.PutShellcode(data); err != nil {
		return err
	}

	return nil
}

func (s *internalService) DeleteShellcode(id int) error {
	if govalidator.IsNull(strconv.Itoa(id)) {
		return errors.New("invalid shellcode id")
	}

	return s.repository.DeleteShellcode(id)
}
