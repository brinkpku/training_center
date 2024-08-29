package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

// server
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许所有的跨域请求，生产环境中应进行安全控制
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 处理 WebSocket 连接的函数
func wsHandler(w http.ResponseWriter, r *http.Request) {
	// 将 HTTP 连接升级为 WebSocket 连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("[SERVER]Failed to upgrade connection:", err)
		return
	}
	defer conn.Close()

	for {
		// 读取消息
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				break
			}
			log.Printf("[SERVER]Unexpected close error: %v", err)
			return
		}
		log.Printf("[SERVER]Received: %s", message)
		if messageType == websocket.CloseMessage {
			log.Println("[SERVER]Received close message from client.")
			break
		}

		// 回送消息
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println("[SERVER]Write error:", err)
			break
		}
	}
}

func runWsServer() {
	// 创建一个简单的 HTTP 服务器
	http.HandleFunc("/ws", wsHandler)
	fmt.Println("WebSocket server started at :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	go runWsServer()
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	log.Printf("Connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Dial error:", err)
	}
	defer conn.Close()

	// 向服务器发送消息
	err = conn.WriteMessage(websocket.TextMessage, []byte("Hello, WebSocket!"))
	if err != nil {
		log.Fatal("Write error:", err)
	}
	// 接收服务器的响应
	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Fatal("Read error:", err)
	}
	log.Printf("Received: %s", message)

	err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Fatal("Write error:", err)
	}
	time.Sleep(time.Second)
	os.Exit(0)
}
