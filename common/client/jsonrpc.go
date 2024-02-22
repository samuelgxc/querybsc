package client

import (
	"bytes"
	"fmt"
	rpcjson "github.com/gorilla/rpc/json"
	"net/http"
	"time"
)

//JsonRpc客户端
//import rpcjson "github.com/gorilla/rpc/json"
func CallJsonRpc(url string, method string, param interface{}, reply interface{}) error {
	//组装JsonRpc请求信息
	body, err := rpcjson.EncodeClientRequest(method, param)
	if err != nil {
		return fmt.Errorf("组装JsonRpc请求信息失败:%s", err.Error())
	}

	//创建HttpRequest
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("创建HttpRequest失败:%s", err.Error())
	}
	//req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Type", "application/json")

	//创建HttpClient并发起请求
	client := http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   true, //true:不同HTTP请求之间TCP连接的重用将被阻止（http1.1默认为长连接，此处改为短连接）
			MaxIdleConnsPerHost: 512,  //控制每个主机下的最大闲置连接数目
		},
		Timeout: time.Minute * 10, //Client请求的时间限制,该超时限制包括连接时间、重定向和读取response body时间;Timeout为零值表示不设置超时
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//解析JsonRpc响应信息
	err = rpcjson.DecodeClientResponse(resp.Body, reply)
	if err != nil {
		return err
	}
	return nil
}
