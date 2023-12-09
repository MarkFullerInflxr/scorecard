package conf

import (
	"os"
)

func Routes(s *Server) error {
	API_BASE := os.Getenv("API_BASE")

	v1 := s.Engine.Group(API_BASE) // global root base path, configure for gin
	{
		pg := v1.Group("/fs")
		{
			pg.GET("/sheet/:tableName", s.Spreadsheet.GetSheet())
			pg.POST("/update/:tableName", s.Spreadsheet.UpdateCell())
			pg.POST("/update/:tableName/addRow", s.Spreadsheet.AddRow())
			pg.GET("/ws", s.Spreadsheet.HandleLiveComm())
		}
	}
	err := s.Engine.Run(":80")

	return err
}