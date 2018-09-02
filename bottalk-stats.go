package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/urfave/cli"
)

func getStatsCommands() []cli.Command {

	return []cli.Command{{

		Name:  "stats",
		Usage: "show bottalk stats",
		Action: func(c *cli.Context) error {

			response, err := http.Get("https://oko.bottalk.de/okostat.json")
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
				fmt.Println("Bottalk Current Stats")
				fmt.Println("---")
				for key, val := range m {
					fmt.Println(key+":", val)
				}
			}

			return nil
		},
	}}
}
