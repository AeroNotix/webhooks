package webhooks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type User struct {
	Email        string `json:"email"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Skype        string `json:"skype"`
	LinkedIn     string `json:"linkedin"`
	Twitter      string `json:"twitter"`
	ProjectLimit int64  `json:"projects_limit"`
	ExternalUID  string `json:"external_UID"`
	Provider     string `json:"provider"`
	Bio          string `json:"bio"`
}

func CreateUser(conf ConfigFile, u User) error {
	jsonstr, err := json.Marshal(u)
	if err != nil {
		return err
	}
	resp, err := http.Post(
		fmt.Sprintf("%s/%s?private_token=%s", conf.Endpoint, "users", conf.APIKey),
		"application/json", strings.NewReader(string(jsonstr)),
	)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}
