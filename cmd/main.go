package main

import (
	"flag"
	"fmt"
	"nacosleak/models"
	"nacosleak/utils"
	"os"
)


func main(){
	flag.StringVar(&utils.Url,"url","","url")
	flag.StringVar(&utils.Usrname,"u","","username")
	flag.StringVar(&utils.Passwd,"p","","password")
	flag.Parse()
	fmt.Println(utils.Banner())
	if utils.Url != ""{
		namespace := models.GetNameSpace(utils.Url)
		fmt.Println(namespace)
		err := models.GetConfig(utils.Url,namespace)
		if err != nil {
			os.Exit(-1)
		}
		utils.SaveConfig(utils.Url,namespace)
		models.Analyze()
	}
}
