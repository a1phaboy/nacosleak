package main

import (
	"flag"
	"fmt"
	"nacosleak/models"
	"nacosleak/utils"
	"os"
)

var targets []string

func main(){
	flag.StringVar(&utils.Url,"t","","目标nacos的url")
	flag.StringVar(&utils.UrlsFile,"ts","","批量扫描，含有URLs的txt文件路径")
	flag.StringVar(&utils.Usrname,"u","","用户名（可选）")
	flag.StringVar(&utils.Passwd,"p","","密码（可选）")
	flag.StringVar(&utils.BasePath,"s","","保存到指定路径（可选）")
	flag.Parse()
	fmt.Println(utils.Banner())
	if utils.UrlsFile != ""{
		urls,err := utils.ReadTargetFile(utils.UrlsFile)
		if err != nil {
			fmt.Println("打开文件失败")
			return
		}
		targets = append(targets,urls...)
	}else if utils.Url != "" {
		targets = append(targets,utils.Url)
	}
	fmt.Println("[*] 目标Total:",len(targets))
	for _,target := range targets{
		utils.Nginx = false
		fmt.Println("[+] start target: ",target)
		namespace,err := models.GetNameSpace(target)
		if err != nil {
			fmt.Println("[-] Get ",target," namespace fail.")
			continue
		}
		fmt.Println("NameSpace: ")
		for _,v := range namespace{
			fmt.Println("   ",v.Name)
		}
		err = models.GetConfig(target,namespace)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}
		utils.SaveConfig(target,namespace)
		models.Analyze()

	}
}
