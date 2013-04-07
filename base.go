package webhooks

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type ConfigFile struct {
	GitURL           string
	Endpoint         string
	APIKey           string
	Username         string
	PublicallyListed string
}

func baseRequest(path, method string, body io.Reader) ([]byte, int, error) {
	req, err := http.NewRequest(
		method,
		path,
		body,
	)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Accept", "application/json")
	if err != nil {
		return nil, 0, err
	}
	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, 0, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	switch resp.StatusCode {
	case http.StatusOK,
		http.StatusCreated:
		return b, resp.StatusCode, nil
	default:
		type Message struct {
			M string `json:"message"`
		}
		m := &Message{}
		err = json.Unmarshal(b, m)
		if err != nil {
			return nil, resp.StatusCode, err
		}
		return nil, resp.StatusCode, errors.New(m.M)
	}
	panic("Unreachable: base.baseRequest")
}

func ReadConfigFile(path string) (ConfigFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return ConfigFile{}, err
	}
	body, err := ioutil.ReadAll(f)
	if err != nil {
		return ConfigFile{}, err
	}
	conf := ConfigFile{}
	err = json.Unmarshal(body, &conf)
	return conf, err
}

func ReadConfigFileFromHome() (ConfigFile, error) {
	home := os.Getenv("HOME")
	if home == "" {
		return ConfigFile{}, errors.New("Cannot determine $HOME variable.")
	}
	newconf, err := ReadConfigFile(filepath.Join(home, ".gitlabclirc"))
	if err != nil {
		return ConfigFile{}, err
	}
	return newconf, nil
}
