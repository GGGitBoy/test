package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var urlMap map[string]bool //防止重复链接进入死循环，不过获取链接太多可能会内存溢出

func fetch(url string, count int) {
	if count > 1 { //设定爬取深度是1页
		return
	}

	body, err := goquery.NewDocument(url)
	if err != nil {
		return
	}
	body.Find("#__docusaurus > div > div > div > div > div > ul > li").Each(func(i int, aa *goquery.Selection) {
		fmt.Printf("%v", aa)
		href, IsExist := aa.Attr("href")
		if IsExist == true {
			href = strings.TrimSpace(href)
			if strings.Contains(href, "docs/rancher2") {
				if len(href) > 2 && IsUrl(href) {
					if _, ok := urlMap[href]; ok == false {
						if strings.HasPrefix(href, "/") || strings.HasPrefix(href, "./") {
							href = SamePathUrl(url, href, 1)
						} else if strings.HasPrefix(href, "../") {
							href = SamePathUrl(url, href, 2)
						}

						fmt.Println("修改之后的url：", href)
						urlMap[href] = true

						fetch(href, count+1)
					}
				}
			}
		}
	})
}

func writeValues(outfile string) error {
	file, err := os.Create(outfile)
	if err != nil {
		fmt.Printf("创建%s文件失败！", outfile)
		return err
	}
	defer file.Close()
	for k, _ := range urlMap {
		file.WriteString(k + "\n")
	}
	return nil
}

func main() {
	urlMap = make(map[string]bool, 1000000)
	fetch("https://docs.rancher.cn/docs/rancher2/releases/v2.4.8/", 0)
	writeValues("urls.dat")
}

//////////

func IsUrl(str string) bool {
	if strings.HasPrefix(str, "#") || strings.HasPrefix(str, "//") || strings.HasSuffix(str, ".exe") || strings.HasSuffix(str, ":void(0);") {
		return false
	} else if strings.HasPrefix(str, "{") && strings.HasSuffix(str, "}") {
		return false
	} else if strings.EqualFold(str, "javascript:;") {
		return false
	} else {
		return true
	}
	return true
}

func SamePathUrl(preUrl string, url string, mark int) (newUrl string) {
	last := strings.LastIndex(preUrl, "/")
	if last == 6 {
		newUrl = preUrl + url
	} else {
		if mark == 1 {
			newUrl = preUrl[:last] + url
		} else {
			newPreUrl := preUrl[:last]
			newLast := strings.LastIndex(newPreUrl, "/")
			newUrl = newPreUrl[:newLast] + url
		}
	}
	return newUrl
}
