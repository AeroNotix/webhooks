package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"webhooks"
)

var conf webhooks.ConfigFile

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

type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Commit struct {
	ID        string `json:"id"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	URL       string `json:"url"`
	Author    Author `json:"author"`
}

type WebHookCommit struct {
	Before     string `json:"before"`
	After      string `json:"after"`
	Ref        string `json:"ref"`
	UserID     int64  `json:"user_id"`
	UserName   string `json:"user_name"`
	Repository struct {
		Name        string `json:"name"`
		URL         string `json:"url"`
		Description string `json:"description"`
		Homepage    string `json:"homepage"`
	} `json:"repository"`
	Commits           []Commit `json:"commits"`
	TotalCommitsCount int64    `json:"total_commits_count"`
}

type SystemHook struct {
	CreatedAt     string `json:"created_at"`
	EventName     string `json:"event_name"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	OwnerEmail    string `json:"owner_email"`
	OwnerName     string `json:"owner_name"`
	UserEmail     string `json:"user_email"`
	Username      string `json:"user_name"`
	Path          string `json:"path"`
	ID            int64  `json:"project_id"`
	ProjectAccess string `json:"project_access"`
	ProjectName   string `json:"project_name"`
	ProjectPath   string `json:"project_path"`
}

func (s SystemHook) NewUserID() (int64, error) {
	users, err := webhooks.ListUsers(conf)
	if err != nil {
		return -1, err
	}
	for _, user := range users {
		if user.Email == s.Email {
			return user.ID, nil
		}
	}
	return -1, errors.New(fmt.Sprintf("User: %s not found", s.Email))
}

func SystemHookEndpoint(w http.ResponseWriter, req *http.Request) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	go func() {
		<-time.After(time.Second * 10)
		hook := &SystemHook{}
		err = json.Unmarshal(b, hook)
		if err != nil {
			log.Println(err)
			return
		}
		switch hook.EventName {
		case "user_create":
			if id, err := hook.NewUserID(); err == nil {
				err := webhooks.AddUserToAllProjects(conf, id, webhooks.GUEST)
				if err != nil {
					log.Println("Error adding user to all projects: " + err.Error())
				}
			} else {
				log.Println("Error retrieving user ID: " + err.Error())
			}
		case "user_add_to_team":
			if hook.Username == "PublicLister" {
				err := webhooks.AddAllUsersToProject(conf, hook.ID, webhooks.GUEST)
				if err != nil {
					log.Println("Error adding all users to project: " + err.Error())
				}
			}
		}
	}()
}

func CommitHandler(w http.ResponseWriter, req *http.Request) {
	if req.Body != nil {
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Println(err)
			return
		}
		whp := &WebHookCommit{}
		err = json.Unmarshal(b, whp)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(whp)
	}
	w.WriteHeader(http.StatusOK)
}

/* Main function */
func main() {
	http.HandleFunc("/", SystemHookEndpoint)
	http.HandleFunc("/commits/", CommitHandler)
	fmt.Println(http.ListenAndServe(":6060", nil))
}
