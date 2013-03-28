package main

import (
	"15.185.120.66/AeroNotix/gitlabcli"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var conf gitlabcli.ConfigFile
var Create = flag.String("create", "", "The name of a repository to create.")
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
		crr, err := gitlabcli.CreateRepository(conf, *Create, nil)
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Println(crr)
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
		crr, err := gitlabcli.CreateRepository(conf, *Init, nil)
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
}
