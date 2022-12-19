package main

import (
	"flag"
	"fmt"
	"nacosleak/analyze"
	"nacosleak/utils"
)


func main(){
	flag.StringVar(&utils.Url,"url","","url")
	flag.StringVar(&utils.Usrname,"u","","username")
	flag.StringVar(&utils.Passwd,"p","","password")
	flag.Parse()
	fmt.Println(utils.Banner())
	if utils.Url != ""{
		config := utils.GetConfig(utils.Url)
		utils.SaveConfig(utils.Url,config)
		analyze.Analyze()
	}
}
