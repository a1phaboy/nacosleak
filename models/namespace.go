package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"nacosleak/utils"
)

func GetNameSpace(url string) (nameSpaces []utils.NacosConfig){
	resp := utils.GetResp(utils.UrlFormat(url)+utils.NAMESPACE_API,false)
	if resp == nil {
		fmt.Println("[ ERROR ] Get namespace fail .")
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("[ ERROR ] read response body fail .")
		return nil
	}
	var bd map[string]interface{}
	json.Unmarshal(body,&bd)
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
