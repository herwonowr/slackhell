package repository

import (
	"errors"

	"go.etcd.io/bbolt"
)

const (
	accountBucket = "account"
)

const (
	_ AccountRole = iota

	// AdminRole ...
	AdminRole

	// AgentRole ...
	AgentRole
)

type (
	// AccountRole ...
	AccountRole int
)

// Account ...
type Account struct {
	ID            int
	Role          AccountRole
	SlackID       string
	SlackRealName string
}

func (s *dbService) initAccount() error {
	err := s.createBucket(accountBucket)
	if err != nil {
		return err
	}

	return nil
}

func (s *dbService) GetAccountByID(id int) (*Account, error) {
	var account Account
	identifier := s.itob(id)

	err := s.getObject(accountBucket, identifier, &account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (s *dbService) GetAccountBySlackID(slackID string) (*Account, error) {
	var account *Account

	err := s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(accountBucket))
		cursor := bucket.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var a Account
			err := s.unmarshalObject(v, &a)
			if err != nil {
				return err
			}

			if a.SlackID == slackID {
				account = &a
				break
			}
		}

		if account == nil {
			return errors.New("account not found")
		}
		return nil
	})

	return account, err
}

func (s *dbService) GetAccountBySlackRealName(slackRealName string) (*Account, error) {
	var account *Account

	err := s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(accountBucket))
		cursor := bucket.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var a Account
			err := s.unmarshalObject(v, &a)
			if err != nil {
				return err
			}

			if a.SlackRealName == slackRealName {
				account = &a
				break
			}
		}

		if account == nil {
			return errors.New("account not found")
		}
		return nil
	})

	return account, err
}

func (s *dbService) GetAccounts() ([]Account, error) {
	var accounts = make([]Account, 0)

	err := s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(accountBucket))
		cursor := bucket.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var a Account
			err := s.unmarshalObject(v, &a)
			if err != nil {
				return err
			}

			accounts = append(accounts, a)
		}

		return nil
	})

	return accounts, err
}

func (s *dbService) CreateAccount(account *Account) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(accountBucket))

		id, _ := bucket.NextSequence()
		account.ID = int(id)

		data, err := s.marshalObject(account)
		if err != nil {
			return err
		}

		return bucket.Put(s.itob(int(account.ID)), data)
	})
}

func (s *dbService) PutAccount(account *Account) error {
	identifier := s.itob(account.ID)
	return s.updateObject(accountBucket, identifier, account)
}

func (s *dbService) DeleteAccount(id int) error {
	identifier := s.itob(id)
	return s.deleteObject(accountBucket, identifier)
}
