package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"webhooks"
)

var conf webhooks.ConfigFile
var ListUsers = flag.Bool("listusers", false, "Lists all users.")
var ListProjects = flag.Bool("listprojects", false, "Lists all projects.")
var Create = flag.String("create", "", "The name of a repository to create.")
var User = flag.Bool("adduser", false, "Add a new user.")
var DelUser = flag.Bool("deluser", false, "Removes a user.")
var DelUserId = flag.Int64("userid", -1, "ID of user to remove.")
var Email = flag.String("email", "", "E-mail address for new user")
var Name = flag.String("name", "", "Name for new user")
var Username = flag.String("username", "", "Username for new user")
var Password = flag.String("password", "", "Password for new user")
var Skype = flag.String("skype", "", "Skype ID for a new user")
var LinkedIn = flag.String("linkedin", "", "LinkedIn ID for a new user")
var Twitter = flag.String("Twitter", "", "Twitter ID for a new user")
var ProjectLimit = flag.Int64("projectlimit", 10, "Project Limit for a new user")
var Init = flag.String("init", "", "The name of a repository to initialize.")

func init() {
	if c, err := webhooks.ReadConfigFileFromHome(); err != nil {
		log.Panic(err)
	} else {
		conf = c
	}
}
func main() {
	flag.Parse()
	if *Create != "" {
		crr, err := webhooks.CreateRepository(conf, *Create, nil)
		if err != nil {
			fmt.Println("Error creating repository: " + err.Error())
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
				return
			}
		}
		user := webhooks.User{
			Email:        *Email,
			Name:         *Name,
			Username:     *Username,
			Password:     *Password,
			Skype:        *Skype,
			LinkedIn:     *LinkedIn,
			Twitter:      *Twitter,
			ProjectLimit: *ProjectLimit,
		}
		fmt.Println(webhooks.CreateUser(conf, user))
		return
	}

	if *DelUser && *DelUserId != -1 {
		fmt.Println(webhooks.DeleteUser(conf, *DelUserId))
		return
	}

	if *ListUsers {
		users, err := webhooks.ListUsers(conf)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, user := range users {
			fmt.Println(fmt.Sprintf(
				`ID: %d
Email: %s
Name: %s
Username: %s
-------------------`, user.ID, user.Email, user.Name, user.Username,
			))
		}
	}

	if *ListProjects {
		projects, err := webhooks.ListProjects(conf)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, project := range projects {
			desc := "<empty>"
			if project.Description != nil {
				desc = *project.Description
			}
			fmt.Println(fmt.Sprintf(
				`ID: %d
Name: %s
Description: %s
Owner: %s
Path: %s
-------------------`,
				project.ID, project.Name, desc, project.Owner.Username, project.Path,
			))
		}
	}
}
