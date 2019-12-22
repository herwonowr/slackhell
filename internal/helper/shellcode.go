package helper

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"strings"
)

// WriteShellcode ...
func (s *Service) WriteShellcode(key string) (string, error) {
	template, err := ioutil.ReadFile("./shellcode/shellcode.code")
	if err != nil {
		return "", err
	}

	hash := md5.New()
	_, err = hash.Write([]byte(key))
	if err != nil {
		return "", err
	}

	hashKey := hex.EncodeToString(hash.Sum(nil))
	shellcode := strings.Replace(string(template), "shellcode_key", hashKey, -1)

	return shellcode, nil
}
