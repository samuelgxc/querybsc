package client

import (
	"encoding/json"
	"fmt"
	"testing"
)

//定义jsonRpc方法入参（参数首字必须大写，客户端与服务端必须一致）
type ParamTest struct {
	Id   int
	Name string
	Age  json.Number
}

//定义jsonRpc方法出参（参数首字必须大写，客户端与服务端必须一致）
type ReplyTest struct {
	Id   int
	Name string
	Age  json.Number
}

//测试JsonRpc客户端
func TestCallJsonRpc_TestStr(t *testing.T) {
	url := "http://127.0.0.1:8083/rpc"
	method := "JsonRpc.TestStr" //与服务端注册的 应用结构名.结构方法名 对应
	param := ParamTest{
		Id:   123,
		Name: "abc",
		Age:  "18",
	}
	var reply string

	err := CallJsonRpc(url, method, param, &reply)
	if err != nil {
		fmt.Printf("jsonRpc请求失败，err:%s\n", err.Error())
	}

	fmt.Println(reply)

	result := ReplyTest{}
	json.Unmarshal([]byte(reply), &result)
	fmt.Println(result.Id)
	fmt.Println(result.Name)
	fmt.Println(result.Age)
}

//测试JsonRpc客户端
func TestCallJsonRpc_TestStruct(t *testing.T) {
	url := "http://127.0.0.1:8083/rpc"
	method := "JsonRpc.TestStruct" //与服务端注册的 应用结构名.结构方法名 对应
	param := ParamTest{
		Id:   123,
		Name: "abc",
		Age:  "18",
	}
	reply := ReplyTest{}

	err := CallJsonRpc(url, method, param, &reply)
	if err != nil {
		fmt.Printf("jsonRpc请求失败，err:%s\n", err.Error())
	}

	fmt.Println(reply)
	fmt.Println(reply.Id)
	fmt.Println(reply.Name)
	fmt.Println(reply.Age)
}
