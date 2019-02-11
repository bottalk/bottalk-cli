package main

import (
	"bytes"
	"encoding/json"
	"github.com/olekukonko/tablewriter"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const apiUrl = "https://api.bottalk.de/"

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

type ResponseSkillList struct {
	ResponseBasic
	Skills []struct {
		Name     string
		Token    string
		Language string
	}
}

type ResponseSkillFiles struct {
	ResponseBasic
	Skill struct {
		Name     string
		Token    string
		Language string
	}
	Files []struct {
		Name    string
		Content string
	}
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

	if strings.Contains(basic.Message, "token expired") {
		log.Println("token expired, trying to refresh")
		getTokenRefresh()
		return apiGet(action, v)
	}

	if len(basic.Error) > 0 {
		log.Fatal("Failed to make a call: " + basic.GetError() + " => " + basic.Message)
		return ""
	}

	return string([]byte(body))
}

func getInfo() {

	info := ResponseInfo{}
	apiGet("info", &info)
	log.Println("You are logged in with account: " + info.User)
}

func getSkillList() {

	info := ResponseSkillList{}
	apiGet("skills", &info)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Language", "Name"})

	for _, v := range info.Skills {
		table.Append([]string{v.Token, v.Language, v.Name})
	}
	table.Render() // Send output
}

func getSkillFiles(skillToken string) {

	info := ResponseSkillFiles{}
	apiGet("skill/"+skillToken, &info)

	path := info.Skill.Name

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0700)
	}

	for _, j := range info.Files {
		ioutil.WriteFile(path+"/"+j.Name, []byte(j.Content), 0700)
	}
}
