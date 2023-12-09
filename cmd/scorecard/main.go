package main

import (
	log "github.com/sirupsen/logrus"
	"influxer/scorecard/cmd/scorecard/conf"
	spreadsheet "influxer/scorecard/pkg/spreadsheet"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	s, routesInit := conf.NewServer()

	s.Spreadsheet = spreadsheet.NewScorecard()

	// build roots
	return routesInit()
}
