package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"nacosleak/utils"
)

func GetNameSpace(url string) (nameSpaces []utils.NacosConfig,err error){
	resp,err,code := utils.GetResp(utils.UrlFormat(url)+utils.NAMESPACE_API,false)
	if err != nil && code == 404{
		resp,err,code = utils.GetResp(utils.UrlFormat(url)+utils.NAMESPACE_API_NGINX,false)
		if err != nil {
			return
		}else{
			utils.Nginx = true
		}
	}
	if resp.StatusCode != 200 {
		fmt.Println("[ ERROR ] Get namespace fail .")
		err = fmt.Errorf("[ ERROR ] Get namespace fail .")
		return nil,err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("[ ERROR ] read response body fail .")
		err = fmt.Errorf("[ ERROR ] read response body fail .")
		return nil,err
	}
	var bd map[string]interface{}
	err = json.Unmarshal(body, &bd)
	if err != nil {
		err = fmt.Errorf("[ ERROR ] handle json.Unmarshal fail.")
		return nil, err
	}
	data := bd["data"].([]interface{})
	for _,v:= range data {
		vv := v.(map[string]interface{})
		var n utils.NacosConfig
		n.Tenant = vv["namespace"].(string)
		n.Name = vv["namespaceShowName"].(string)
		nameSpaces = append(nameSpaces,n)
	}
	return
}
