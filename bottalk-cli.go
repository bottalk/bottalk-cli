package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

//Version of the cli
const Version = 19

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

	checkUpdate()

	err := bottalkCli.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
