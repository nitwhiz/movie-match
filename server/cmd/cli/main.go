package main

import (
	"github.com/nitwhiz/movie-match/server/internal/command"
	"github.com/nitwhiz/movie-match/server/internal/config"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
	})

	if err := config.Init(); err != nil {
		log.Fatal(err)
	}

	if config.C.Debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("running in debug mode")
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

func main() {
	if err := command.GetApp().Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
