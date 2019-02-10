package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const apiUrl = "http://localhost:9098/"

type ResponseBasic struct {
	Result  string
	Error   string
	Message string
}

func (h ResponseBasic) GetError() string {
	return h.Error
}

type ResponseInfo struct {
	ResponseBasic
	User string
}

func apiGet(action string, v interface{}) string {

	url := apiUrl + action

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + token.Token

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)
	req.Header.Set("User-Agent", "bottalk-cli v"+strconv.Itoa(Version))

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error on response.\n[ERRO] -", err)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	jd := json.NewDecoder(bytes.NewReader([]byte(body)))
	err = jd.Decode(&v)
	if err != nil {
		log.Fatal("Failed to decode response: " + err.Error())
		return ""
	}

	basic := ResponseBasic{}
	jd = json.NewDecoder(bytes.NewReader([]byte(body)))
	err = jd.Decode(&basic)
	if err != nil {
		log.Fatal("Failed to decode response: " + err.Error())
		return ""
	}

	if len(basic.Error) > 0 {
		log.Fatal("Failed to make a call: " + basic.GetError() + " => " + basic.Message)
		return ""
	}

	if strings.Contains(basic.Message, "token expired") {
		log.Println("token expired, trying to refresh")
		getTokenRefresh()
		return apiGet(action, v)
	}

	return string([]byte(body))
}

func getInfo() {

	info := ResponseInfo{}
	apiGet("info", &info)
	log.Println("You are logged in with account: " + info.User)
}
