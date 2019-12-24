package helper

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/asaskevich/govalidator"
)

const (
	phpShell = "./shellcode/php/shellcode.code"
	aspShell = "./shellcode/asp/shellcode.code"
)

// WriteShellcode ...
func (s *Service) WriteShellcode(key string, shell string) (string, error) {
	if govalidator.IsNull(key) {
		return "", errors.New("invalid shellcode key")
	}

	if !(shell == "php" || shell == "asp") {
		return "", errors.New("invalid shellcode type, valid shellcode php, asp")
	}

	shellType := shell
	if shell == "php" {
		shellType = phpShell
	}

	if shell == "asp" {
		shellType = aspShell
	}

	template, err := ioutil.ReadFile(shellType)
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
