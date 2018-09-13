package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

//Version of the cli
const Version = 19

var path string

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

	bottalkCli.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "dir",
			Value:       "./",
			Usage:       "Directory to work in",
			Destination: &path,
			EnvVar:      "APP_DIR,DIR",
		},
	}

	checkUpdate()

	err := bottalkCli.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
