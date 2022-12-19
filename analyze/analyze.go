package analyze

import (
	"bufio"
	"fmt"
	"io"
	"nacosleak/utils"
	"os"
	"regexp"
	"strings"
)

func Analyze() {
	var passwdText []string
	fi, err := os.Open(utils.FolderName+"/all_config.txt")
	if err != nil {
		fmt.Println("[ ERROR ]open file fail .")
		os.Exit(0)
	}
	passwdReg := regexp.MustCompile(`(?i)password:`)
	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		res := passwdReg.FindAllStringSubmatch(string(a),-1)
		if res != nil {
			as := string(a)
			index := strings.Index(as,":")
			passwdText = append(passwdText,strings.TrimSpace(as[index+1:]))
		}
	}
	utils.SavePasswd(passwdText)
	//fmt.Println(passwdText)
}
