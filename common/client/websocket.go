package client

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/url"
)

//创建WebSocket客户端连接
//import "github.com/gorilla/websocket"
func NewWebSocketConn(scheme, host, path string) (*websocket.Conn, error) {
	//处理scheme
	if scheme != "ws" && scheme != "wss" {
		scheme = "ws"
	}
	//处理host
	if host == "" {
		return nil, fmt.Errorf("host不能为空")
	}
	//组装url
	Url := url.URL{Scheme: scheme, Host: host, Path: path}
	//创建webSocket连接
	ws, _, err := websocket.DefaultDialer.Dial(Url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("创建webSocket连接出错，url:%s，err: %v", Url.String(), err)
	}
	return ws, nil
}

////创建WebSocket客户端连接
////import "golang.org/x/net/websocket"
//func NewWebSocketHttpConn(scheme, host, path string) (*xwebsocket.Conn, error) {
//	//处理scheme
//	hscheme := "http"
//	if scheme != "ws" && scheme != "wss" {
//		scheme = "ws"
//	} else if scheme == "wss" {
//		hscheme = "https"
//	}
//	//处理host
//	if host == "" {
//		return nil, fmt.Errorf("host不能为空")
//	}
//	//组装url
//	Url := url.URL{Scheme: scheme, Host: host, Path: path}
//	Origin := url.URL{Scheme: hscheme, Host: host, Path: ""}
//	//创建webSocket连接
//	ws, err := xwebsocket.Dial(Url.String(), "", Origin.String())
//	if err != nil {
//		return nil, fmt.Errorf("创建webSocket连接出错，url:%s，err: %v", Url.String(), err)
//	}
//	return ws, nil
//}
