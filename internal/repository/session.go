package repository

import (
	"errors"

	"go.etcd.io/bbolt"
)

const (
	sessionBucket = "session"
)

// Session ...
type Session struct {
	ID  string
	Key string
}

func (s *dbService) initSession() error {
	err := s.createBucket(sessionBucket)
	if err != nil {
		return err
	}

	return nil
}

func (s *dbService) GetSessionByID(id string) (*Session, error) {
	var sesssion *Session

	err := s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(sessionBucket))
		cursor := bucket.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var c Session
			err := s.unmarshalObject(v, &c)
			if err != nil {
				return err
			}

			if c.ID == id {
				sesssion = &c
				break
			}
		}

		if sesssion == nil {
			return errors.New("key not found")
		}
		return nil
	})

	return sesssion, err
}

func (s *dbService) CreateSession(session *Session) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(sessionBucket))

		data, err := s.marshalObject(session)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(session.ID), data)
	})
}

func (s *dbService) PutSession(session *Session) error {
	identifier := []byte(session.ID)
	return s.updateObject(sessionBucket, identifier, session)
}

func (s *dbService) DeleteSession(id string) error {
	identifier := []byte(id)
	return s.deleteObject(sessionBucket, identifier)
}
