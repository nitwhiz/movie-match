package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nitwhiz/movie-match/server/internal/config"
	"github.com/nitwhiz/movie-match/server/internal/dbutils"
	"github.com/nitwhiz/movie-match/server/internal/handler"
	"log"
)

func main() {
	if err := config.Init(); err != nil {
		panic(err)
	}

	db, err := dbutils.GetConnection()

	if err != nil {
		log.Fatal(err)
	}

	if err := dbutils.Migrate(db); err != nil {
		log.Fatal(err)
	}

	if err := dbutils.InitUsers(db); err != nil {
		log.Fatal(err)
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

	err = router.Run("0.0.0.0:6445")

	if err != nil {
		panic(err)
	}
}
