package main

import (
	"15.185.120.66/AeroNotix/webhooks"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var conf webhooks.ConfigFile
var Create = flag.String("create", "", "The name of a repository to create.")
var User = flag.Bool("adduser", false, "Add a new user")
var Email = flag.String("email", "", "E-mail address for new user")
var Username = flag.String("username", "", "Username for new user")
var Password = flag.String("password", "", "Password for new user")
var Skype = flag.String("skype", "", "Skype ID for a new user")
var LinkedIn = flag.String("linkedin", "", "LinkedIn ID for a new user")
var Twitter = flag.String("Twitter", "", "Twitter ID for a new user")
var ProjectLimit = flag.Int64("projectlimit", 10, "Project Limit for a new user")
var Init = flag.String("init", "", "The name of a repository to initialize.")

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
func main() {
	flag.Parse()
	if *Create != "" {
		crr, err := webhooks.CreateRepository(conf, *Create, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("New repository created with ID: %d\n", crr.ID)
		return
	}
	if *Init != "" {
		/* Run the Git command to create a new repository locally */
		cmd := exec.Command("git", "init", *Init)
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		crr, err := webhooks.CreateRepository(conf, *Init, nil)
		if err != nil {
			log.Fatal(err)
		}
		/* Move into the new sub directory */
		err = os.Chdir(*Init)
		if err != nil {
			log.Fatal(err)
		}
		/* Add the remote as the origin of the new repository */
		cmd = exec.Command(
			"git", "remote", "add", "origin", conf.GitURL+crr.PathWithNS,
		)
		cmd.Stdout = os.Stdout
		err = cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	if *User {
		for field, errmsg := range map[string]string{
			*Email:    "Missing e-mail.",
			*Username: "Missing username.",
			*Password: "Missing password.",
		} {
			if field == "" {
				fmt.Println(errmsg)
				goto after_user_create
			}
		}
		user := webhooks.User{
			Email:        *Email,
			Username:     *Username,
			Password:     *Password,
			Skype:        *Skype,
			LinkedIn:     *LinkedIn,
			Twitter:      *Twitter,
			ProjectLimit: *ProjectLimit,
		}
		fmt.Println(webhooks.CreateUser(conf, user))
	}
after_user_create:
}
