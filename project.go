package gitlabcli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type owner struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	State     string `json:"state"`
	CreatedAt string `json:"created_at"`
}
type namespace struct {
	CreatedAt   string  `json:"created_at"`
	Description *string `json:"description"`
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	OwnerID     int64   `json:"owner_id"`
	Path        string  `json:"path"`
	UpdatedAt   string  `json:"updated_at"`
}

type CreateRepositoryResponse struct {
	ID                   int64     `json:"id"`
	Description          *string   `json:"description"`
	DefaultBranch        *string   `json:"default_branch"`
	Owner                owner     `json:"owner"`
	Public               bool      `json:"public"`
	Path                 string    `json:"path"`
	PathWithNS           string    `json:"path_with_namespace"`
	IssuesEnabled        bool      `json:"issues_enabled"`
	MergeRequestsEnabled bool      `json:"merge_requests_enabled"`
	WallEnabled          bool      `json:"wall_enabled"`
	WikiEnabled          bool      `json:"wiki_enabled"`
	CreatedAt            string    `json:"created_at"`
	Namespace            namespace `json:"namespace"`
}

/* Will create a new repository using the global conf object. */
func CreateRepository(conf ConfigFile, r string, extra *map[string]string) (*CreateRepositoryResponse, error) {
	path := fmt.Sprintf(conf.Endpoint+"projects?name=%s&", r)
	vals := url.Values{}
	if extra != nil {
		for k, v := range *extra {
			vals.Add(k, v)
		}
	}
	path = path + vals.Encode()
	req, err := http.NewRequest("POST",
		path,
		nil,
	)
	c := http.Client{}
	if err != nil {
		return nil, err
	}
	req.Header.Add("PRIVATE-TOKEN", conf.APIKey)
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case http.StatusCreated:
		crr := &CreateRepositoryResponse{}
		err = json.Unmarshal(body, crr)
		if err != nil {
			return nil, err
		}
		return crr, nil
	default:
		type Message struct {
			M string `json:"message"`
		}
		var msg Message
		err = json.Unmarshal(body, &msg)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(msg.M)
	}
	panic("Unreachable!")

}
