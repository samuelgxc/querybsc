package client

import (
	"fmt"
	"testing"
)

//测试Http客户端
func TestCallHttp(t *testing.T) {
	method := "POST"
	url := "http://127.0.0.1:8082/test/json?id=123"
	params := map[string]interface{}{}
	params["name"] = "abc"
	params["account"] = 200.005
	headers := map[string]string{}
	headers["Content-Type"] = "application/json"

	//发起 http json 请求
	reply, err := CallHttp(method, url, params, headers)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(string(reply))
}
