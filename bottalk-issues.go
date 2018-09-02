package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/urfave/cli"
)

func getIssueCommands() []cli.Command {
	return []cli.Command{
		{
			Category: "Issue tracker",
			Name:     "list",
			Usage:    "list bottalk issues",
			Aliases:  []string{"ls"},
			Action: func(c *cli.Context) error {

				response, err := http.Get("https://dev.bottalk.de/api/issue")
				if err != nil {
					fmt.Printf("%s", err)
					os.Exit(1)
				} else {
					defer response.Body.Close()
					contents, err := ioutil.ReadAll(response.Body)
					if err != nil {
						fmt.Printf("%s", err)
						os.Exit(1)
					}

					m := map[string]interface{}{}
					json.Unmarshal([]byte(contents), &m)
					fmt.Println("Bottalk Active Issues")
					fmt.Println("---")
					for key, val := range m {
						fmt.Println(key+":", val)
					}
				}

				return nil
			},
		},
		{
			Category: "Issue tracker",
			Name:     "add",
			Aliases:  []string{"a"},
			Usage:    "add a task to the list",
			Action: func(c *cli.Context) error {
				fmt.Println("Adding task: ", c.Args().First())

				values := map[string]string{"name": c.Args().First()}

				jsonValue, _ := json.Marshal(values)

				response, err := http.Post("https://dev.bottalk.de/api/issue", "application/json", bytes.NewBuffer(jsonValue))
				if err != nil {
					fmt.Printf("%s", err)
					os.Exit(1)
				} else {
					defer response.Body.Close()
					body, _ := ioutil.ReadAll(response.Body)
					m := map[string]interface{}{}
					json.Unmarshal([]byte(body), &m)
					if m["result"] == "ok" {
						fmt.Println("Successfully added task with id", m["id"])
					} else {
						fmt.Println("Couldn't add task: ", m["error"])
					}
				}
				return nil
			},
		},
		{
			Category: "Issue tracker",
			Name:     "complete",
			Aliases:  []string{"c"},
			Usage:    "complete a task in the list",
			Action: func(c *cli.Context) error {
				fmt.Println("Finishing task: ", c.Args().First())

				values := map[string]string{"id": c.Args().First()}

				jsonValue, _ := json.Marshal(values)

				response, err := http.Post("https://dev.bottalk.de/api/completeissue", "application/json", bytes.NewBuffer(jsonValue))
				if err != nil {
					fmt.Printf("%s", err)
					os.Exit(1)
				} else {
					defer response.Body.Close()
					body, _ := ioutil.ReadAll(response.Body)
					m := map[string]interface{}{}
					json.Unmarshal([]byte(body), &m)
					if m["result"] == "ok" {
						fmt.Println("Successfully completed task:", m["name"])
					} else {
						fmt.Println("Couldn't complete task: ", m["error"])
					}
				}
				return nil
			},
		},
	}
}
