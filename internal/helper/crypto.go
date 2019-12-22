package helper

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

// GenerateRandomByte ...
func (s *Service) GenerateRandomByte(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString ...
func (s *Service) GenerateRandomString(n int) (string, error) {
	b, err := s.GenerateRandomByte(n)
	return base64.URLEncoding.EncodeToString(b), err
}
