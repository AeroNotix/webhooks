package webhooks

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
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
