package webhooks

import (
	"encoding/json"
	"fmt"
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
	body, _, err := baseRequest(
		fmt.Sprintf("%s/%s?private_token=%s", conf.Endpoint, "users", conf.APIKey),
		"GET",
		nil,
	)
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
	_, statuscode, err := baseRequest(
		fmt.Sprintf("%s/%s?private_token=%s", conf.Endpoint, "users", conf.APIKey),
		"POST",
		strings.NewReader(string(jsonstr)),
	)
	if err != nil || statuscode != http.StatusCreated {
		return err
	}
	return nil
}

func DeleteUser(conf ConfigFile, ID int64) error {
	_, statuscode, err := baseRequest(
		"DELETE",
		fmt.Sprintf("%s/users/%d?private_token=%s", conf.Endpoint, ID, conf.APIKey),
		nil,
	)
	if err != nil || statuscode != http.StatusOK {
		return err
	}
	return nil
}

func ListUsersForProject(conf ConfigFile, ID int64) ([]User, error) {
	body, statuscode, err := baseRequest(
		fmt.Sprintf("%s/projects/%d/members?private_token=%s", conf.Endpoint, ID, conf.APIKey),
		"GET",
		nil,
	)
	if err != nil || statuscode != http.StatusOK {
		return nil, err
	}
	users := []User{}
	err = json.Unmarshal(body, &users)
	return users, err
}
