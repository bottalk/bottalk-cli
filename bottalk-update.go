package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func checkUpdate() {
	//log.Println("Checking version...")
	response, err := http.Get("https://oko.bottalk.de/console/version.json")
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		m := map[string]interface{}{}
		json.Unmarshal([]byte(body), &m)
		d := m["version"]
		if d.(float64) > float64(Version) {
			log.Println("Got new version of app, performing binary update")

			updatedBinary, err := http.Get("https://oko.bottalk.de/console/binary." +
				fmt.Sprintf("%v", m["version"].(float64)) + ".bin")
			if err != nil {
				fmt.Printf("%s", err)
				os.Exit(1)
			}

			if updatedBinary.StatusCode != 200 {
				log.Println("Got wrong status while updating, skipping op")
				return
			}

			ex, err := os.Executable()
			if err != nil {
				panic(err)
			}

			log.Println("Replacing: ", ex)

			body, _ := ioutil.ReadAll(updatedBinary.Body)
			ioutil.WriteFile(ex+".new", body, 0777)

			os.Rename(ex+".new", ex)

			log.Println("Successfully wrote out new version. Please restart this command.")
			os.Exit(2)

		}
	}
}
