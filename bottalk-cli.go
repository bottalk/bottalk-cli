package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

//Version of the cli
const Version = 19

var path string

type Token struct {
	Token   string `json:"access_token"`
	Expires int    `json:"expires_in"`
	Refresh string `json:"refresh_token"`
	Scope   string `json:"scope"`
}

var token Token

func main() {
	bottalkCli := cli.NewApp()
	bottalkCli.Name = "bottalk"
	bottalkCli.Usage = "CLI helper for Bottalk.de"
	bottalkCli.Description = ""
	bottalkCli.Version = fmt.Sprintf("%d", Version)

	bottalkCli.Commands = append(
		bottalkCli.Commands,
		getTestCommands()...,
	)

	bottalkCli.Commands = append(
		bottalkCli.Commands,
		getLoginCommands()...,
	)

	bottalkCli.Commands = append(
		bottalkCli.Commands,
		getSkillCommands()...,
	)

	bottalkCli.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "dir",
			Value:       "./",
			Usage:       "Directory to work in",
			Destination: &path,
			EnvVar:      "APP_DIR,DIR",
		},
	}

	loadToken()

	checkUpdate()

	err := bottalkCli.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func loadToken() bool {
	f, _ := os.Open(userHomeDir() + "/.bottalk-cli.token")
	j := json.NewDecoder(f)
	err := j.Decode(&token)
	if err != nil {
		log.Println("Wrong data in tokenfile: " + err.Error())
		return false
	}
	return true
}
