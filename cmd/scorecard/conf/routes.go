package conf

import (
	"os"
)

func Routes(s *Server) error {
	API_BASE := os.Getenv("API_BASE")

	v1 := s.Engine.Group(API_BASE) // global root base path, configure for gin
	{
		favi := v1.Group("/public")
		{
			favi.GET("/:file", s.Favicon.Handle())
		}
		pg := v1.Group("/fs")
		{
			pg.GET("/sheet/", s.Spreadsheet.ListSheets())
			pg.GET("/sheet/:tableName", s.Spreadsheet.GetSheet())
			pg.GET("/live/:tableName", s.Spreadsheet.HandleLiveComm())
			pg.GET("/live/~", s.Spreadsheet.HandleLiveComm())
		}
	}
	err := s.Engine.Run(":80")

	return err
}
