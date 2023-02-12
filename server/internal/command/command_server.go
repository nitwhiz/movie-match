package command

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nitwhiz/movie-match/server/internal/dbutils"
	"github.com/nitwhiz/movie-match/server/internal/handler"
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

	router := gin.Default()

	router.Use(cors.Default())

	handler.AddMediaRetrieveAll(router, db)
	handler.AddMediaRetrieveById(router, db)
	handler.AddMediaRetrievePoster(router, db)

	handler.AddUserRetrieveAll(router, db)

	handler.AddUserMediaAddSeen(router, db)
	handler.AddUserMediaVote(router, db)

	handler.AddUserMatchRetrieveAll(router, db)

	handler.AddUserMediaRecommendedRetrieveAll(router, db)

	return router.Run("0.0.0.0:6445")
}
