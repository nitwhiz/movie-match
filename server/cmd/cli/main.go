package main

import (
	"github.com/nitwhiz/movie-match/server/internal/command"
	"github.com/nitwhiz/movie-match/server/internal/config"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	if err := config.Init(); err != nil {
		log.Fatal(err)
	}

	app := command.GetApp()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
