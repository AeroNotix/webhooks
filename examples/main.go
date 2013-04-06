package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

type SystemHookProjects struct {
	CreatedAt  string `json:"created_at"`
	EventName  string `json:"event_name"`
	Name       string `json:"name"`
	OwnerEMail string `json:"owner_email"`
	OwnerName  string `json:"owner_name"`
	Path       string `json:"path"`
	ID         int64  `json:"project_id"`
}

type SystemHookUsers struct {
	CreatedAt  string `json:"created_at"`
	EventName  string `json:"event_name"`
	Name       string `json:"name"`
	OwnerEMail string `json:"email"`
}

func Users(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Users")
	if req.Body != nil {
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Println(err)
			return
		}
		whu := &SystemHookUsers{}
		err = json.Unmarshal(b, whu)
		if err != nil {
			log.Println(err)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func Projects(w http.ResponseWriter, req *http.Request) {
	if req.Body != nil {
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Println(err)
			return
		}
		whcp := &SystemHookProjects{}
		err = json.Unmarshal(b, whcp)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(whcp)
	}
	w.WriteHeader(http.StatusOK)
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
	http.HandleFunc("/users/", Users)
	http.HandleFunc("/projects/", Projects)
	fmt.Println(http.ListenAndServe(":6060", nil))
}
