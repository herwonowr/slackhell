package service

import (
	"errors"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/herwonowr/slackhell/internal/repository"
)

func (s *internalService) GetAccountByID(id int) (*repository.Account, error) {
	if govalidator.IsNull(strconv.Itoa(id)) {
		return nil, errors.New("invalid account id")
	}

	return s.repository.GetAccountByID(id)
}

func (s *internalService) GetAccountBySlackID(slackID string) (*repository.Account, error) {
	if govalidator.IsNull(slackID) {
		return nil, errors.New("invalid account Slack id")
	}

	return s.repository.GetAccountBySlackID(slackID)
}

func (s *internalService) GetAccountBySlackRealName(slackRealName string) (*repository.Account, error) {
	if govalidator.IsNull(slackRealName) {
		return nil, errors.New("invalid account slack real name")
	}

	return s.repository.GetAccountBySlackRealName(slackRealName)
}

func (s *internalService) GetAccounts() ([]repository.Account, error) {
	return s.repository.GetAccounts()
}

func (s *internalService) CreateAccount(role repository.AccountRole, slackID string, slackRealName string) error {
	if govalidator.IsNull(slackID) {
		return errors.New("invalid account slack id")
	}

	if govalidator.IsNull(slackRealName) {
		return errors.New("invalid account slack real name")
	}

	aSlackID, err := s.repository.GetAccountBySlackID(slackID)
	if err != nil && err.Error() != "account not found" {
		return err
	}

	if aSlackID != nil {
		return errors.New("account slack id already exists")
	}

	aSlackRealName, err := s.repository.GetAccountBySlackRealName(slackRealName)
	if err != nil && err.Error() != "account not found" {
		return err
	}

	if aSlackRealName != nil {
		return errors.New("account slack real name already exists")
	}

	vRole := repository.AgentRole
	if !(role == repository.AdminRole || role == repository.AgentRole) {
		return errors.New("invalid account role")
	}

	if role == repository.AdminRole {
		vRole = repository.AdminRole
	}

	if role == repository.AgentRole {
		vRole = repository.AgentRole
	}

	account := &repository.Account{
		Role:          vRole,
		SlackID:       slackID,
		SlackRealName: slackRealName,
	}

	err = s.repository.CreateAccount(account)
	if err != nil {
		return err
	}

	return nil
}
func (s *internalService) PutAccount(id int, role repository.AccountRole, slackID string, slackRealName string) error {
	if govalidator.IsNull(strconv.Itoa(id)) {
		return errors.New("invalid account id")
	}

	account, err := s.repository.GetAccountByID(id)
	if err != nil {
		return err
	}

	vRole := account.Role
	vSlackID := account.SlackID
	vSlackRealName := account.SlackRealName

	if role != 0 {
		if !(role == repository.AdminRole || role == repository.AgentRole) {
			return errors.New("invalid account role")
		}

		if role == repository.AdminRole {
			vRole = repository.AdminRole
		}

		if role == repository.AgentRole {
			vRole = repository.AgentRole
		}
	}

	if !govalidator.IsNull(slackID) {
		aSlackID, err := s.repository.GetAccountBySlackID(slackID)
		if err != nil {
			return err
		}

		if aSlackID != nil && aSlackID.ID != account.ID {
			return errors.New("slack id already exists")
		}

		vSlackID = slackID
	}

	if !govalidator.IsNull(slackRealName) {
		aSlackRealName, err := s.repository.GetAccountBySlackRealName(slackRealName)
		if err != nil {
			return err
		}

		if aSlackRealName != nil && aSlackRealName.ID != account.ID {
			return errors.New("account already exists")
		}

		vSlackRealName = slackRealName
	}

	data := &repository.Account{
		ID:            account.ID,
		Role:          vRole,
		SlackID:       vSlackID,
		SlackRealName: vSlackRealName,
	}

	if err := s.repository.PutAccount(data); err != nil {
		return err
	}

	return nil
}

func (s *internalService) DeleteAccount(id int) error {
	if govalidator.IsNull(strconv.Itoa(id)) {
		return errors.New("invalid account id")
	}

	return s.repository.DeleteAccount(id)
}
