package sideload

import (
	"github.com/nitwhiz/movie-match/server/internal/config"
	log "github.com/sirupsen/logrus"
	"os"
)

func IsConfigured() bool {
	return config.C.Sideload != nil
}

func HasPostersFsAccess() bool {
	if _, err := os.Stat(config.C.Sideload.Posters.FsBasePath); err != nil {
		log.Error(err)

		return false
	}

	return true
}
