package client

import (
	"fmt"
	"github.com/gorilla/websocket"
	"strings"
	"testing"
	"time"
)

//测试WebSocket客户端连接
//import "github.com/gorilla/websocket"
func TestNewWebSocketConn(t *testing.T) {
	scheme := "ws"
	host := "127.0.0.1:8081"
	path := "/ws"
	ws, err := NewWebSocketConn(scheme, host, path)
	if err != nil {
		fmt.Println("连接ws失败", err)
		return
	}
	defer ws.Close()

	sendMsg := make(chan string, 10000)

	//读取ws信息
	go func(ws *websocket.Conn) {
		defer func() {
			ws.Close()
		}()
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				fmt.Println("读取ws失败", err)
				return
			}
			msg := string(message)
			fmt.Println("读取信息", msg)

			//读取到ping信息，发送pong信息
			if strings.Contains(msg, "ping") {
				msg := strings.Replace(msg, "ping", "pong", 1)
				sendMsg <- msg
			}
		}
	}(ws)

	//发送ws信息
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		//定时发送ping信息
		case ti := <-ticker.C:
			msg := fmt.Sprintf(`{"ping":%v}`, ti.Unix())
			sendMsg <- msg
			//发送信息
		case msg := <-sendMsg:
			fmt.Println("发送信息", msg)
			if err := ws.WriteMessage(websocket.BinaryMessage, []byte(msg)); err != nil {
				fmt.Println("发送信息失败", err)
				return
			}
		}
	}
}

////测试WebSocket客户端连接
////import "golang.org/x/net/websocket"
//func TestNewWebSocketHttpConn(t *testing.T) {
//	scheme := "ws"
//	host := "127.0.0.1:8084"
//	path := "/ws"
//	ws,err := NewWebSocketHttpConn(scheme, host, path)
//	if err != nil {
//		fmt.Println("连接ws失败", err)
//		return
//	}
//	defer ws.Close()
//
//	sendMsg := make(chan string, 10000)
//
//	//读取ws信息
//	go func(ws *xwebsocket.Conn) {
//		defer func() {
//			ws.Close()
//		}()
//		for {
//			var message []byte
//			err := xwebsocket.Message.Receive(ws, &message)
//			if err != nil {
//				fmt.Println("读取ws失败", err)
//				return
//			}
//			msg := string(message)
//			fmt.Println("读取信息", msg)
//
//			//读取到ping信息，发送pong信息
//			if strings.Contains(msg, "ping") {
//				msg := strings.Replace(msg, "ping", "pong", 1)
//				sendMsg <- msg
//			}
//		}
//	}(ws)
//
//	//发送ws信息
//	ticker := time.NewTicker(time.Second * 5)
//	for {
//		select {
//		//定时发送ping信息
//		case ti := <-ticker.C:
//			msg := fmt.Sprintf(`{"ping":%v}`, ti.Unix())
//			sendMsg <- msg
//		//发送信息
//		case msg := <-sendMsg:
//			fmt.Println("发送信息", msg)
//			if err := xwebsocket.Message.Send(ws, []byte(msg)); err != nil {
//				fmt.Println("发送信息失败", err)
//				return
//			}
//		}
//	}
//}
