package service

import "github.com/herwonowr/slackhell/internal/repository"

// Service ...
type Service interface {
	GetAccountByID(id int) (*repository.Account, error)
	GetAccountBySlackID(slackID string) (*repository.Account, error)
	GetAccountBySlackRealName(slackRealName string) (*repository.Account, error)
	GetAccounts() ([]repository.Account, error)
	CreateAccount(role repository.AccountRole, slackID string, slackRealName string) error
	PutAccount(id int, role repository.AccountRole, slackID string, slackRealName string) error
	DeleteAccount(id int) error

	GetShellcodeByID(id int) (*repository.Shellcode, error)
	GetShellcodeByKey(key string) (*repository.Shellcode, error)
	GetShellcodes() ([]repository.Shellcode, error)
	CreateShellcode(fileID string, shellType string, shellkey string, ownerID string, ownerRealName string) error
	PutShellcode(key string, endpoint string, ownerID string) error
	DeleteShellcode(id int) error

	GetSessionByID(id string) (*repository.Session, error)
	CreateSession(id string, key string) error
	PutSession(id string, key string) error
	DeleteSession(id string) error

	GetVersion() (int, error)
	PutVersion(version int) error
}
type internalService struct {
	repository repository.Service
}

// NewService ...
func NewService(repo repository.Service) Service {
	return &internalService{
		repository: repo,
	}
}
