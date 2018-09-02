package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

type filedata struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func getTestCommands() []cli.Command {

	return []cli.Command{{

		Name:  "test",
		Usage: "Start test in the current folder",
		Action: func(c *cli.Context) error {
			log.Println("Gathering files in current folder")
			matches, _ := filepath.Glob("*.yml")
			if len(matches) == 0 {
				log.Println("Scenario files not found -- nothing to test")
				os.Exit(1)
			}

			sendFiles := []filedata{}

			for _, file := range matches {
				log.Println("* " + file)
				dat, _ := ioutil.ReadFile(file)
				sendFiles = append(sendFiles, filedata{Name: file, Content: string(dat)})
			}

			values := map[string][]filedata{"files": sendFiles}

			jsonValue, _ := json.Marshal(values)

			response, err := http.Post("https://dev.bottalk.de/api/externaltest", "application/json", bytes.NewBuffer(jsonValue))

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
				log.Println("Sending skill files to bottalk")
				log.Println("---")
				if m["result"] != "ok" {
					for key, val := range m {
						log.Println("* "+key+":", val)
					}
					os.Exit(1)
				} else {
					log.Println("Temprorary skill created successfully")
					os.Exit(0)
				}
			}
			return nil
		},
	}}
}
