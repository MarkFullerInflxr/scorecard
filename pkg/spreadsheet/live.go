package spreadsheet

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	utils "influxer/scorecard/utilities"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (s *Scorecard) HandleLiveComm() gin.HandlerFunc {
	return func(c *gin.Context) {
		tableName := c.Param("tableName")
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		s.conn[tableName] = append(s.conn[tableName], conn)
		returnTable := s.BuildTable(tableName)
		err = conn.WriteMessage(websocket.TextMessage, returnTable)
		if err != nil {
			return
		}

		go s.iDidntKnowYouKnewHowToRead(conn)
	}
}

func (s *Scorecard) iDidntKnowYouKnewHowToRead(comm *websocket.Conn) {
	for {
		messageType, buf, err := comm.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		event := make(map[string]interface{})
		fmt.Println(string(buf))
		err = utils.FromJsonBytes(buf, &event)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(messageType, event)

		s.handleEvent(event, comm)
	}
}
