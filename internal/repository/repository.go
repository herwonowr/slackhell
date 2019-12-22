package repository

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"time"

	"go.etcd.io/bbolt"
)

// Service ...
type Service interface {
	GetAccountByID(id int) (*Account, error)
	GetAccountBySlackID(slackID string) (*Account, error)
	GetAccountBySlackRealName(slackRealName string) (*Account, error)
	GetAccounts() ([]Account, error)
	CreateAccount(account *Account) error
	PutAccount(account *Account) error
	DeleteAccount(id int) error

	GetShellcodeByID(id int) (*Shellcode, error)
	GetShellcodeByKey(key string) (*Shellcode, error)
	GetShellcodes() ([]Shellcode, error)
	CreateShellcode(shellcode *Shellcode) error
	PutShellcode(shellcode *Shellcode) error
	DeleteShellcode(id int) error

	GetSessionByID(id string) (*Session, error)
	CreateSession(session *Session) error
	PutSession(session *Session) error
	DeleteSession(id string) error

	GetVersion() (int, error)
	PutVersion(version int) error

	Init() error
	Close() error
}

type dbService struct {
	db *bbolt.DB
}

// NewService ...
func NewService(dbPath string) (Service, error) {
	db, err := bbolt.Open(dbPath, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	return &dbService{db}, nil
}

func (s *dbService) Init() error {
	err := s.initVersion()
	if err != nil {
		return err
	}
	err = s.initAccount()
	if err != nil {
		return err
	}
	err = s.initShellcode()
	if err != nil {
		return err
	}
	err = s.initSession()
	if err != nil {
		return err
	}

	return nil
}

func (s *dbService) Close() error {
	return s.db.Close()
}

func (s *dbService) createBucket(bucketName string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *dbService) itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func (s *dbService) getObject(bucketName string, key []byte, object interface{}) error {
	var data []byte

	err := s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))

		value := bucket.Get(key)
		if value == nil {
			return errors.New("object not found")
		}

		data = make([]byte, len(value))
		copy(data, value)

		return nil
	})
	if err != nil {
		return err
	}

	return s.unmarshalObject(data, object)
}

func (s *dbService) updateObject(bucketName string, key []byte, object interface{}) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))

		data, err := s.marshalObject(object)
		if err != nil {
			return err
		}

		err = bucket.Put(key, data)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *dbService) deleteObject(bucketName string, key []byte) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		return bucket.Delete(key)
	})
}

func (s *dbService) getNextIdentifier(bucketName string) int {
	var identifier int

	s.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		id, err := bucket.NextSequence()
		if err != nil {
			return err
		}
		identifier = int(id)
		return nil
	})

	return identifier
}

func (s *dbService) marshalObject(object interface{}) ([]byte, error) {
	return json.Marshal(object)
}

func (s *dbService) unmarshalObject(data []byte, object interface{}) error {
	return json.Unmarshal(data, object)
}
