package utils

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)
const NAMESPACE_API = "/nacos/v1/console/namespaces"
const NAMESPACE_API_NGINX = "/v1/console/namespaces"
const CONFIG_API = "/nacos/v1/cs/configs?dataId=&group=&appName=&config_tags=&pageNo=1&pageSize=100&tenant=%s&search=accurate"
const CONFIG_API_NGINX = "/v1/cs/configs?dataId=&group=&appName=&config_tags=&pageNo=1&pageSize=100&tenant=%s&search=accurate"
const LOGIN_API = "/nacos/v1/auth/users/login"
const LOGIN_API_NGINX = "/v1/auth/users/login"
var Usrname string
var Passwd string
var Url string
var FolderName string
var NameSpaceFolder []string

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
		fmt.Println("[ ERROR ] Cannot connect to target, plz check out ")
		os.Exit(0)
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

func UrlFormat (url string ) string {
	if strings.LastIndex(url, "/") == len(url)-1 {
		url = url[:len(url)-1]
	}
	return url
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

func SaveConfig(url string, configs []NacosConfig) bool {
	domain := url[7:]
	FolderName = strings.Replace(domain, ".", "_", -1)
	FolderName = strings.Replace(FolderName, ":", "_", -1)
	basePath,_ := os.Getwd()
	foldPath := filepath.Join(basePath,"results",FolderName)
	if !exists(foldPath) {
		err := os.MkdirAll(foldPath, 0766)
		if err != nil {
			fmt.Println(err)
		}
	}
	for _,config := range configs{
		nameSpacePath := filepath.Join(foldPath,config.Name)
		if !exists(nameSpacePath) {
			err := os.MkdirAll(nameSpacePath, 0766)
			if err != nil {
				fmt.Println(err)
			}
		}
		//fmt.Println(nameSpacePath)
		NameSpaceFolder = append(NameSpaceFolder,nameSpacePath)
		allConf := nameSpacePath + "/all_config.txt"
		f, err := os.Create(allConf)
		if err != nil {
			fmt.Println("[ ERROR ] create file fail .")
			return false
		} else {
			_, err = f.Write([]byte(config.Config))
		}
	}
	fmt.Println("[ SUCCESS ] Save in path:" + foldPath)
	return true
}

func SavePasswd(path string,passwdz []string) bool {
	passwdText := path + "/passwd.txt"
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

func SaveAKSK(path string,akskz []string) bool {
	akskText := path + "/aksk.txt"
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
                                         Author:a1  v1.2            `
}
