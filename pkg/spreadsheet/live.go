package spreadsheet

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (s *Scorecard) HandleLiveComm() gin.HandlerFunc {
	return func(c *gin.Context) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		s.comm = conn
		err = conn.WriteMessage(websocket.TextMessage, []byte("Hell Yea Brother"))
		if err != nil {
			return
		}

		go iDidntKnowYouKnewHowToRead(conn)
	}
}

func iDidntKnowYouKnewHowToRead(comm *websocket.Conn) {
	for {
		messageType, buf, err := comm.ReadMessage()
		if err != nil {
			return
		}
		fmt.Println(messageType, string(buf))
	}
}
