package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

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
type WebHookPost struct {
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

/* This is a basic handler which shows an example hook firing */
func HandleWebHook(w http.ResponseWriter, req *http.Request) {
	if req.Body != nil {
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		whp := &WebHookPost{}
		err = json.Unmarshal(b, whp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		fmt.Println(whp)
	}
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/", HandleWebHook)
	fmt.Println(http.ListenAndServe(":12346", nil))
}
