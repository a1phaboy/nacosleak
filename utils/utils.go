package utils

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/url"
	"os"
	"strings"
)

func Banner() string {
	return `  
  _   _    _    ____ ___  ____  _     _____    _    _  __
 | \ | |  / \  / ___/ _ \/ ___|| |   | ____|  / \  | |/ /
 |  \| | / _ \| |  | | | \___ \| |   |  _|   / _ \ | ' / 
 | |\  |/ ___ \ |__| |_| |___) | |___| |___ / ___ \| . \ 
 |_| \_/_/   \_\____\___/|____/|_____|_____/_/   \_\_|\_\
                                         By:a1  v2.1            `
}

func ReadTargetFile(path string) (targets []string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		targets = append(targets, scanner.Text())
	}
	return targets, scanner.Err()
}

func UrlFormat(Url string) string {
	var u *url.URL
	if strings.HasPrefix(Url, "http") || strings.HasPrefix(Url, "https") {
		u, _ = url.Parse(Url)
		return fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	} else {
		u, _ = url.Parse("//" + Url)
		u.Scheme = "http"
	}
	return fmt.Sprintf("%s://%s", u.Scheme, u.Host)
}

// 生成一个随机的6位字符串
func GenerateRandomString() string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var sb strings.Builder
	for i := 0; i < 6; i++ {
		// 随机选择letters中的一个字符
		randomChar, err := getRandomChar(letters)
		if err != nil {
			panic(err)
		}
		sb.WriteByte(randomChar)
	}
	return sb.String()
}

// 从给定的字符串中随机选择一个字符
func getRandomChar(letters string) (byte, error) {
	// 生成一个[0, len(letters)-1]范围内的随机数
	randIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
	if err != nil {
		return 0, err
	}
	// 将随机数转换为字符
	return letters[randIndex.Int64()], nil
}
