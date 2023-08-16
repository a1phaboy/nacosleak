package models

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"nacosleak/utils"
	"net/http"
)

func GetConfig(url string,namespace []utils.NacosConfig) (err error) {
	for i,v := range namespace {
		var resp *http.Response
		if !utils.Nginx  {
			if v.Name == "public" {
				resp,err,_ = utils.GetResp(utils.UrlFormat(url)+fmt.Sprintf(utils.CONFIG_API,""), false)
			}else{
				resp,err,_ = utils.GetResp(utils.UrlFormat(url)+fmt.Sprintf(utils.CONFIG_API,v.Tenant), false)
			}
			if err != nil {
				err = fmt.Errorf("[ ERROR ] Get Resp fail .")
				return
			}
		}else{
			if v.Name == "public" {
				resp,err,_ = utils.GetResp(utils.UrlFormat(url)+fmt.Sprintf(utils.CONFIG_API_NGINX,""), false)
			}else{
				resp,err,_ = utils.GetResp(utils.UrlFormat(url)+fmt.Sprintf(utils.CONFIG_API_NGINX,v.Tenant), false)
			}
			if err != nil {
				err = fmt.Errorf("[ ERROR ] Get Resp fail .")
				return
			}
		}

		var body []byte
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			err = fmt.Errorf("[ ERROR ] Get Resp's body fail .")
			return
		}
		Config := make(map[string]interface{})
		if string(body) != ""{
			err = json.Unmarshal(body, &Config)
			if err != nil {
				err = fmt.Errorf("json.Unmarshal fail")
				fmt.Println(body)
				return
			}
		}
		var yaml_data interface{}
		var config string
		if Config["pageItems"] != nil {
			for _, v := range Config["pageItems"].([]interface{}) {
				item := v.(map[string]interface{})
				config = config + fmt.Sprintf("\n--------------------      %v     --------------------\n",item["dataId"].(string))
				yaml.Unmarshal([]byte(item["content"].(string)), &yaml_data)
				config = config + item["content"].(string)
			}
			fmt.Println("[ SUCCESS ] Get configs on nacos success .")
		}
		//return config
		//configs = append(configs,config)
		namespace[i].Config = config
	}
	return

}
