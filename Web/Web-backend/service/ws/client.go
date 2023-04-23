package ws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"mp333player.com/models"
	"mp333player.com/utils/e"
	"net/http"
	"time"
)

type SendMsg struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
}

type ReplyMsg struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

var successMsg = ReplyMsg{
	Code: e.WebsocketSuccessMessage,
	Msg:  "OK!",
}

var errorMsg = ReplyMsg{
	Code: e.ERROR,
	Msg:  "ERROR!",
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Hub *Hub

	// The websocket connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		//先要关闭通道连接
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		//开始接受发送过来的数据【如果接收的数据是json格式，可以用ReadJSON()】
		sendMsg := new(SendMsg)
		err := c.Conn.ReadJSON(&sendMsg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		s, _ := json.Marshal(successMsg)
		e, _ := json.Marshal(errorMsg)
		//1:获取指纹id；2：录入指纹；3：获取温度；4：获取指纹失败；5：获取温度失败
		if sendMsg.Type == 4 {
			//获取指纹失败
			_ = c.Conn.WriteMessage(websocket.TextMessage, e)
		} else if sendMsg.Type == 5 {
			//获取温度失败
			_ = c.Conn.WriteMessage(websocket.TextMessage, e)
		} else if sendMsg.Type == 1 {
			//成功获取到指纹
			fingerprint := models.FingerprintLog{
				FingerId: sendMsg.Content,
				WorkNum:  "不存在",
				Status:   "找不到用户！",
			}
			err = models.DB.Create(&fingerprint).Error
			_ = c.Conn.WriteMessage(websocket.TextMessage, s)
		} else if sendMsg.Type == 2 {
			//录入指纹
			fingerprint := models.FingerprintLog{
				FingerId: sendMsg.Content,
				WorkNum:  "待录入",
				Status:   "新用户",
			}
			_ = models.DB.Create(&fingerprint).Error
			_ = c.Conn.WriteMessage(websocket.TextMessage, s)
		} else if sendMsg.Type == 3 {
			//成功获取到温度数据
			temp := models.TemperatureLog{
				TemperatureData: sendMsg.Content,
			}
			_ = models.DB.Create(&temp).Error
			_ = c.Conn.WriteMessage(websocket.TextMessage, s)
		} else {
			_ = c.Conn.WriteMessage(websocket.TextMessage, s)
		}

		//处理字节流,并广播
		//message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		//c.hub.broadcast <- message

	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		//读取通道里的消息
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{Hub: hub, Conn: conn, Send: make(chan []byte, 256)}
	client.Hub.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
