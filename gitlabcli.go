package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
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

func init() {
	home := os.Getenv("HOME")
	if home == "" {
		panic("Cannot determine $HOME variable!")
	}
	f, err := os.Open(filepath.Join(home, ".gitlabclirc"))
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &conf)
	if err != nil {
		panic(err)
	}
}

/* Will create a new repository using the global conf object. */
func CreateRepository(r string, extra *map[string]string) (*CreateRepositoryResponse, error) {
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

type ConfigFile struct {
	Endpoint string
	APIKey   string
}

var conf ConfigFile
var Create = flag.String("create", "", "The name of a repository to create.")
var Init = flag.String("init", "", "The name of a repository to initialize.")

func main() {
	flag.Parse()
	if *Create != "" {
		crr, err := CreateRepository(*Create, nil)
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Println(crr)
		return
	}
	if *Init != "" {
		cmd := exec.Command("git", "init", *Init)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		crr, err := CreateRepository(*Init, nil)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(crr)
		return
	}
}
