package command

import (
	"github.com/nitwhiz/movie-match/server/internal/auth"
	"github.com/nitwhiz/movie-match/server/internal/controller"
	"github.com/nitwhiz/movie-match/server/internal/dbutils"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func Server(_ *cli.Context) error {
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

	router, err := controller.Init(db)

	if err != nil {
		log.Error("Router Init Error: ", err)
		return err
	}

	(auth.NewTokenCleanup(db)).Start()

	return router.Run("0.0.0.0:6445")
}
