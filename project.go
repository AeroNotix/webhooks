package webhooks

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
)

type AccessLevel int

const (
	GUEST     = 10
	REPORTER  = 20
	DEVELOPER = 30
	MASTER    = 40
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
	path := fmt.Sprintf("%s/projects?private_token=%s", conf.Endpoint, conf.APIKey)
	body, statuscode, err := baseRequest(
		path,
		"GET",
		nil,
	)
	if err != nil || statuscode != http.StatusOK {
		return nil, err
	}
	projects := []Project{}
	err = json.Unmarshal(body, &projects)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

/* Will create a new repository using the global conf object. */
func CreateRepository(conf ConfigFile, r string, extra *map[string]string) (*Project, error) {
	path := fmt.Sprintf("%s/projects?private_token=%s&projects&name=%s&", conf.Endpoint, conf.APIKey, r)
	vals := url.Values{}
	if extra != nil {
		for k, v := range *extra {
			vals.Add(k, v)
		}
	}
	path = path + vals.Encode()
	body, statuscode, err := baseRequest(
		path,
		"POST",
		nil,
	)
	if err != nil || statuscode != http.StatusCreated {
		return nil, err
	}
	project := &Project{}
	err = json.Unmarshal(body, project)
	return project, err
}

func AddAllUsersToAllProjects() error {
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

func AddUserToAllProjects(conf ConfigFile, ID int64, a AccessLevel) error {
	oldkey := conf.PublicallyListed
	conf.APIKey = conf.PublicallyListed
	projects, err := ListProjects(conf)
	if err != nil {
		return err
	}
	conf.APIKey = oldkey
	for _, project := range projects {
		if project.IsPublicallyListed(conf) {
			err = project.AddUser(conf, ID, a)
			if err != nil {
				log.Println(err)
			}
		}
	}
	return nil
}

func AddAllUsersToProject(conf ConfigFile, ID int64, a AccessLevel) error {
	return nil
}

func (p Project) AddUser(conf ConfigFile, ID int64, a AccessLevel) error {
	_, statuscode, err := baseRequest(
		fmt.Sprintf(
			"%s/projects/%d/members?private_token=%s&user_id=%d&access_level=%d",
			conf.Endpoint, p.ID, conf.APIKey, ID, a,
		),
		"POST",
		nil,
	)
	if err != nil || statuscode != http.StatusOK {
		return err
	}
	return nil
}

func (p Project) IsPublicallyListed(conf ConfigFile) bool {
	users, err := ListUsersForProject(conf, p.ID)
	if err != nil {
		log.Println(err)
		return false
	}
	for _, user := range users {
		if user.Name == "PublicallyListed" {
			return true
		}
	}
	return false
}
