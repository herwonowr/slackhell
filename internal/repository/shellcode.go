package repository

import (
	"errors"

	"go.etcd.io/bbolt"
)

const (
	shellcodeBucket = "shellcode"
)

// Shellcode ...
type Shellcode struct {
	ID            int
	FileID        string
	ShellKey      string
	Endpoint      string
	OwnerID       string
	OwnerRealName string
}

func (s *dbService) initShellcode() error {
	err := s.createBucket(shellcodeBucket)
	if err != nil {
		return err
	}

	return nil
}

func (s *dbService) GetShellcodeByID(id int) (*Shellcode, error) {
	var shellcode Shellcode
	identifier := s.itob(id)

	err := s.getObject(shellcodeBucket, identifier, &shellcode)
	if err != nil {
		return nil, err
	}

	return &shellcode, nil
}

func (s *dbService) GetShellcodeByKey(key string) (*Shellcode, error) {
	var shellcode *Shellcode

	err := s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(shellcodeBucket))
		cursor := bucket.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var c Shellcode
			err := s.unmarshalObject(v, &c)
			if err != nil {
				return err
			}

			if c.ShellKey == key {
				shellcode = &c
				break
			}
		}

		if shellcode == nil {
			return errors.New("key not found")
		}
		return nil
	})

	return shellcode, err
}

func (s *dbService) GetShellcodes() ([]Shellcode, error) {
	var shellcodes = make([]Shellcode, 0)

	err := s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(shellcodeBucket))
		cursor := bucket.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var c Shellcode
			err := s.unmarshalObject(v, &c)
			if err != nil {
				return err
			}

			shellcodes = append(shellcodes, c)
		}

		return nil
	})

	return shellcodes, err
}

func (s *dbService) CreateShellcode(shellcode *Shellcode) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(shellcodeBucket))

		id, _ := bucket.NextSequence()
		shellcode.ID = int(id)

		data, err := s.marshalObject(shellcode)
		if err != nil {
			return err
		}

		return bucket.Put(s.itob(int(shellcode.ID)), data)
	})
}

func (s *dbService) PutShellcode(shellcode *Shellcode) error {
	identifier := s.itob(shellcode.ID)
	return s.updateObject(shellcodeBucket, identifier, shellcode)
}

func (s *dbService) DeleteShellcode(id int) error {
	identifier := s.itob(id)
	return s.deleteObject(shellcodeBucket, identifier)
}
