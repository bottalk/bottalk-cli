package main

import (
	"bytes"
	"encoding/json"
	"github.com/olekukonko/tablewriter"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	pathlib "path"
	"path/filepath"
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

type ResponseCreateSkill struct {
	ResponseBasic
	Token string
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

func apiPost(action string, data []byte, v interface{}) string {

	url := apiUrl + action

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + token.Token

	// Create a new request using http
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))

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
		return apiPost(action, data, v)
	}

	if len(basic.Error) > 0 {
		log.Fatal("Failed to make a call: " + basic.GetError() + " => " + basic.Message)
		return ""
	}

	return string([]byte(body))
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

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(map[string]string{"token": info.Skill.Token, "name": info.Skill.Name, "language": info.Skill.Language})

	ioutil.WriteFile(path+"/"+".skill.manifest", b.Bytes(), 0700)
	log.Println("Successfully wrote skill " + info.Skill.Token + " (" + info.Skill.Name + ")")
}

func createNewSkill(skillName string, skillLanguage string) {

	if skillLanguage == "" {
		skillLanguage = "en-US"
	}
	values := map[string]string{"name": skillName, "language": skillLanguage}

	jsonValue, _ := json.Marshal(values)

	postResponse := ResponseCreateSkill{}

	apiPost("skills/new", jsonValue, &postResponse)

	if postResponse.Result == "ok" {

		path := skillName

		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.Mkdir(path, 0700)
		}

		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(map[string]string{"token": postResponse.Token, "name": skillName, "language": skillLanguage})

		ioutil.WriteFile(path+"/"+".skill.manifest", b.Bytes(), 0700)
		log.Println("Successfully created skill '" + postResponse.Token + "' (" + skillName + ")")

		log.Println("Token: " + postResponse.Token)
	}
}

func pushSkillFiles() {

	manifest, err := os.Open(".skill.manifest")

	if err != nil {
		log.Fatalln("Cannot read skill manifest: " + err.Error())
	}

	mf := struct {
		Name     string
		Token    string
		Language string
	}{}

	err = json.NewDecoder(manifest).Decode(&mf)
	if err != nil {
		log.Fatalln("Cannot decode manifest: " + err.Error())
	}

	log.Println("Gathering files in folder: ", path)
	matches, _ := filepath.Glob(path + "*.yml")
	if len(matches) == 0 {
		log.Println("Scenario files not found -- nothing to push")
		os.Exit(0)
	}

	sendFiles := []filedata{}

	for _, file := range matches {
		log.Println("* " + file)
		dat, _ := ioutil.ReadFile(file)
		sendFiles = append(sendFiles, filedata{Name: pathlib.Base(file), Content: string(dat)})
	}

	values := map[string][]filedata{"files": sendFiles}

	jsonValue, _ := json.Marshal(values)

	postResponse := ResponseBasic{}

	apiPost("skill/"+mf.Token, jsonValue, &postResponse)

	if postResponse.Result == "ok" {
		log.Println("Successfully pushed your skill files")
	}

}
