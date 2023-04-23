package ws

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//升级websocket协议
func wshandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	//连接成功，创建一个连接实例
	client := &Client{Hub: hub, Conn: conn, Send: make(chan []byte, 256)}
	//将连接实例放入hub
	client.Hub.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
func HttpController(c *gin.Context, hub *Hub) {
	wshandler(hub, c.Writer, c.Request)
}
