package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bregydoc/gtranslate"
	"github.com/mmcdole/gofeed"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func SendDing(title, desc, link string)  {
	postData := make(map[string]interface{}, 0)
	postData["msgtype"] = "actionCard"
	postData["actionCard"] = map[string]string{
		"title":title,
		"text":desc,
		"singleTitle":"阅读全文",
		"singleURL":link,
	}
	bytePostData,_ := json.Marshal(postData)
	resp, err := http.Post("https://oapi.dingtalk.com/robot/send?access_token=4d4d3709b502d0bb6ec50d2cab1daa55c360737da902dad1151631b887fc2ef5", "application/json", bytes.NewReader(bytePostData))
	if err!=nil {
		log.Println("send ding data:", string(bytePostData))
		log.Println("send ding err:", err.Error())
	}
	respByte,_ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Println(string(respByte))
}



func translateTextG(text string) (string ,error) {

	translated, err := gtranslate.TranslateWithParams(
		text,
		gtranslate.TranslationParams{
			From: "en",
			To:   "zh-cn",
		},
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("en: %s | zh: %s \n", text, translated)
	return translated,nil
}

func QueryTw(uname string) error {
	fp := gofeed.NewParser()
	//file, _ := os.Open("DIYgod.xml")
	//defer file.Close()
	//feed, _ := fp.Parse(file)
	feedUrl := "http://103.80.132.31:1200/twitter/user/"+uname+"/exclude_replies"
	feed, err := fp.ParseURL(feedUrl)
	//fmt.Println(feedUrl)
	if err!=nil {
		log.Println("parse url err:", err.Error())
		return err
	}
	if feed==nil {
		log.Println("feed nil")
		return errors.New("feed nil")
	}
	//fmt.Println(feed.Title)
	//fmt.Println(feed.Updated)


	cacheFileName := "./cache/"+uname+"_lastQuery.md"
	_, err = os.Stat(cacheFileName)
	var lasted int64 = 0
	if os.IsNotExist(err) {
		f,err := os.Create(cacheFileName)
		if err!=nil {
			fmt.Println("创建文件或者打开文件失败")
			log.Fatalln(err.Error())
		}
		f.Close()
	}else{
		lastByte,err := os.ReadFile(cacheFileName)
		if err!=nil {
			fmt.Println("打开文件失败")
			log.Fatalln(err.Error())
		}
		lasted,_ = strconv.ParseInt(string(lastByte), 10, 64)
	}

	for _,item := range feed.Items {
		if lasted>= item.PublishedParsed.Unix(){
			//log.Println("未到更新时间")
			continue
		}

		fmt.Println("lasted:", item.PublishedParsed)
		fmt.Println("title:", item.Title)
		fmt.Println()
		lasted = item.PublishedParsed.Unix()
		ioutil.WriteFile(cacheFileName, []byte(strconv.FormatInt(lasted, 10)), 0644)
		if item.Description=="" {
			continue
		}
		item.Description = strings.ReplaceAll(item.Description, "<br>", "\n\n")

		reg := regexp.MustCompile(`<(\S*?)[^>]*>.*?|<.*? />`)
		item.Description = reg.ReplaceAllString(item.Description, "")

		fyStr, err := translateTextG(item.Description)
		if err != nil {
			log.Fatalln("err fy", err.Error())
		}
		fmt.Println(fyStr)
		desc :="来自:"+uname+ "\n\n"+ "翻译:"+fyStr+"\n\n"+"原文:"+item.Description
		SendDing(item.Title, desc, item.Link)
		//resp,err := http.PostForm("http://www.58meishi.cn:8999", url.Values{"desc":{desc}})
		//if err!=nil {
		//	log.Println("err:", err.Error())
		//}else{
		//	respb,_ := ioutil.ReadAll(resp.Body)
		//	fmt.Println(string(respb))
		//}
		// 等待3s 继续下一个
		time.Sleep(3*time.Second)
	}


	return nil
}


