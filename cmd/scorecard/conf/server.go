package conf

import (
	"github.com/gin-gonic/gin"
	"influxer/scorecard/cmd/scorecard/middlewares"
	"influxer/scorecard/pkg/spreadsheet"
)

type Server struct {
	Engine      *gin.Engine
	Spreadsheet spreadsheet.Spreadsheet
}

func NewServer() (*Server, func() error) {
	s := Server{}

	// confgure gin
	g := gin.Default()
	g.Use(gin.Recovery())
	s.Engine = g

	g.Use(middlewares.OptionsMiddleware())
	g.Use(middlewares.CORSMiddleware())

	// any other server specific initialization
	// don't initialize the services in here,
	// it'll make it harder to mock later
	return &s, func() error { return Routes(&s) }
}
