package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}

func bAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func getTokenByCode(code string) string {
	var username = "4"
	var passwd = "bottalk-cli"
	client := &http.Client{}

	body := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {code},
		"redirect_uri": {"http://localhost:2876/auth"},
	}

	reqBody := bytes.NewBufferString(body.Encode())

	req, err := http.NewRequest("POST", "https://auth.bottalk.de/token", reqBody)
	req.Header.Add("Authorization", "Basic "+bAuth(username, passwd))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	s := string(bodyText)

	homedir := userHomeDir()

	d1 := []byte(s)
	err = ioutil.WriteFile(homedir+"/.bottalk-cli.token", d1, 0644)

	if err != nil {
		log.Println(err)
	} else {
		log.Println("Successfully wrote credentials.")
	}

	return s
}

func getTokenRefresh() string {
	var username = "4"
	var passwd = "bottalk-cli"
	client := &http.Client{}

	body := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {token.Refresh},
		"scope":         {token.Scope},
	}

	reqBody := bytes.NewBufferString(body.Encode())

	req, err := http.NewRequest("POST", "https://auth.bottalk.de/token", reqBody)
	req.Header.Add("Authorization", "Basic "+bAuth(username, passwd))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	s := string(bodyText)

	homedir := userHomeDir()

	d1 := []byte(s)
	err = ioutil.WriteFile(homedir+"/.bottalk-cli.token", d1, 0644)

	if err != nil {
		log.Println(err)
	} else {
		log.Println("Successfully wrote credentials.")
	}

	loadToken()

	return s
}

func getLoginCommands() []cli.Command {

	s := &http.Server{
		Addr: ":2876",
	}

	return []cli.Command{{

		Name:  "login",
		Usage: "Login into your bottalk account",
		Action: func(c *cli.Context) error {
			log.Println("Please log in into your account in window that was opened in your browser")
			openbrowser("https://auth.bottalk.de/authorize?client_id=4&response_type=code&scope=all&redirect_uri=http%3A%2F%2Flocalhost%3A2876%2Fauth")
			log.Println("If it didn't open, please follow this link: https://auth.bottalk.de/authorize?client_id=4&response_type=code&scope=all&redirect_uri=http%3A%2F%2Flocalhost%3A2876%2Fauth")

			http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
				code := r.FormValue("code")
				fmt.Fprint(w, " <script>setTimeout(function(){window.close();},5000);</script> This window will close in 5 seconds. You can return to your console.")

				go getTokenByCode(code)

				go func() {
					time.Sleep(1 * time.Second)
					s.Shutdown(context.TODO())

				}()

			})
			err := s.ListenAndServe()
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
			return nil
		},
	}, {

		Name:  "info",
		Usage: "Information about your bottalk account",
		Action: func(c *cli.Context) error {
			log.Println("Fetching info about your account")
			getInfo()
			return nil
		},
	}}

}

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	} else if runtime.GOOS == "linux" {
		home := os.Getenv("XDG_CONFIG_HOME")
		if home != "" {
			return home
		}
	}
	return os.Getenv("HOME")
}
