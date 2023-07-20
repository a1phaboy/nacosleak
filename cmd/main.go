package main

import (
	"flag"
	"fmt"
	"nacosleak/models"
	"nacosleak/utils"
	"os"
)


func main(){
	flag.StringVar(&utils.Url,"t","","目标nacos的url")
	flag.StringVar(&utils.Usrname,"u","","用户名（可选）")
	flag.StringVar(&utils.Passwd,"p","","密码（可选）")
	flag.StringVar(&utils.BasePath,"s","","保存到指定路径（可选）")
	flag.Parse()
	fmt.Println(utils.Banner())
	if utils.Url != ""{
		namespace := models.GetNameSpace(utils.Url)
		fmt.Println("NameSpace: ")
		for _,v := range namespace{
			fmt.Println("   ",v.Name)
		}

		err := models.GetConfig(utils.Url,namespace)
		if err != nil {
			os.Exit(-1)
		}
		utils.SaveConfig(utils.Url,namespace)
		models.Analyze()
	}
}
