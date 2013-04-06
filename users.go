package webhooks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type User struct {
	ID           int64   `json:"id"`
	Email        string  `json:"email"`
	Name         string  `json:"name"`
	Blocked      bool    `json:"blocked"`
	CreatedAt    string  `json:"create_at"`
	Username     string  `json:"username"`
	Password     string  `json:"password"`
	Skype        string  `json:"skype"`
	LinkedIn     string  `json:"linkedin"`
	Twitter      string  `json:"twitter"`
	ProjectLimit int64   `json:"projects_limit"`
	ExternalUID  *string `json:"external_UID"`
	Provider     *string `json:"provider"`
	Bio          *string `json:"bio"`
	DarkScheme   bool    `json:"dark_scheme"`
	ThemeID      int64   `json:"theme_id"`
}

func ListUsers(conf ConfigFile) ([]User, error) {
	resp, err := http.Get(
		fmt.Sprintf("%s/%s?private_token=%s", conf.Endpoint, "users", conf.APIKey),
	)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	users := []User{}
	return users, json.Unmarshal(body, &users)
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

func AddUserToAllProjects(conf ConfigFile, ID int64) error {
	for _, project := range ListProjects(conf) {

	}
	return nil
}
