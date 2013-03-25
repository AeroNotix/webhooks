package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

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

func CreateRepository(r string, extra *map[string]string) error {
	req, err := http.NewRequest("POST",
		fmt.Sprintf(conf.Endpoint+"projects?name=%s", r),
		nil,
	)
	c := http.Client{}
	if err != nil {
		return err
	}
	req.Header.Add("PRIVATE-TOKEN", conf.APIKey)
	if extra != nil {
		for k, v := range *extra {
			req.Header.Add(k, v)
		}
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}

type ConfigFile struct {
	Endpoint string
	APIKey   string
}

var conf ConfigFile
var Create = flag.String("create", "", "The string of a repository to create")

func main() {
	flag.Parse()
	if *Create != "" {
		err := CreateRepository(*Create, nil)
		if err != nil {
			fmt.Println(err)
		}
	}
}
