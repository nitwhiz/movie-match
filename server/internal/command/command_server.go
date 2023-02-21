package command

import (
	"github.com/nitwhiz/movie-match/server/internal/dbutils"
	"github.com/nitwhiz/movie-match/server/internal/server"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func Server(context *cli.Context) error {
	db, err := dbutils.GetConnection()

	if err != nil {
		return err
	}

	if err := dbutils.Migrate(db); err != nil {
		return err
	}

	if err := dbutils.InitUsers(db); err != nil {
		return err
	}

	var opts []server.Option

	if context.Bool(FlagServerWeb) {
		log.Info("Web Enabled")
		opts = append(opts, server.WithRouter())
	}

	if context.Bool(FlagServerWithTokenCleanup) {
		log.Info("Token Cleanup Enabled")
		opts = append(opts, server.WithTokenCleanup())
	}

	if context.Bool(FlagServerWithMediaAutoPull) {
		log.Info("Auto Media Pull Enabled")
		opts = append(opts, server.WithAutoPull())
	}

	s, err := server.New(db, opts...)

	log.Info("server initialized")

	if err != nil {
		return err
	}

	return s.Start()
}
