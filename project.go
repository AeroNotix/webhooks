package webhooks

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
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

type Project struct {
	ID                   int64     `json:"id"`
	Name                 string    `json:"name"`
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

func ListProjects(conf ConfigFile) ([]Project, error) {
	path := conf.Endpoint + "projects"
	req, err := http.NewRequest("GET",
		path,
		nil,
	)
	c := http.Client{}
	if err != nil {
		fmt.Println("new req")
		return nil, err
	}
	req.Header.Add("PRIVATE-TOKEN", conf.APIKey)
	resp, err := c.Do(req)
	if err != nil {
		fmt.Println("do")
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(resp.StatusCode)
	switch resp.StatusCode {
	case http.StatusOK:
		projects := []Project{}
		err = json.Unmarshal(body, &projects)
		if err != nil {
			return nil, err
		}
		return projects, nil
	}
	panic("Unreachable!")
}

/* Will create a new repository using the global conf object. */
func CreateRepository(conf ConfigFile, r string, extra *map[string]string) (*Project, error) {
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
		fmt.Println("new req")
		return nil, err
	}
	req.Header.Add("PRIVATE-TOKEN", conf.APIKey)
	resp, err := c.Do(req)
	if err != nil {
		fmt.Println("do")
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case http.StatusCreated:
		crr := &Project{}
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

func AddUsersToAllProjects() error {
	err := os.Chdir("/home/git/gitlab/")
	if err != nil {
		return err
	}
	cmd := exec.Command("bundle", "exec", "rake",
		"gitlab:import:all_users_to_all_projects", "RAILS_ENV=production",
	)
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
