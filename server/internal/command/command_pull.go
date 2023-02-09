package command

import (
	"github.com/nitwhiz/movie-match/server/internal/dbutils"
	"github.com/nitwhiz/movie-match/server/internal/provider"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func PullMedia(context *cli.Context) error {
	providerName := context.Args().Get(0)
	mediaProvider := provider.GetByName(providerName)

	if mediaProvider == nil {
		log.Fatalf("media provider '%s' not found.", providerName)
	}

	if err := mediaProvider.Init(); err != nil {
		log.Fatalf(err.Error())
	}

	db, err := dbutils.GetConnection()

	if err != nil {
		return err
	}

	return mediaProvider.Pull(db)
}
