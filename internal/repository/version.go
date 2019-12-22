package repository

import (
	"errors"
	"strconv"

	"go.etcd.io/bbolt"
)

const (
	versionBucket = "version"
	versionKey    = "VERSION"
)

func (s *dbService) initVersion() error {
	err := s.createBucket(versionBucket)
	if err != nil {
		return err
	}

	return nil
}

func (s *dbService) GetVersion() (int, error) {
	var data []byte

	err := s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(versionBucket))

		value := bucket.Get([]byte(versionKey))
		if value == nil {
			return errors.New("version not found")
		}

		data = make([]byte, len(value))
		copy(data, value)

		return nil
	})
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(string(data))
}

func (s *dbService) PutVersion(version int) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(versionBucket))

		data := []byte(strconv.Itoa(version))
		return bucket.Put([]byte(versionKey), data)
	})
}
