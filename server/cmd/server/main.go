package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nitwhiz/movie-match/server/internal/config"
	"github.com/nitwhiz/movie-match/server/internal/handler"
	"github.com/nitwhiz/movie-match/server/pkg/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func main() {
	if err := config.Init(); err != nil {
		panic(err)
	}

	dsn := "host=localhost user=root password=root dbname=movie_match sslmode=disable TimeZone=Europe/Berlin"

	verboseLogger := logger.New(
		log.New(os.Stdout, "\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Error,
			Colorful:      true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: verboseLogger,
	})

	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(
		&model.User{},
		&model.Media{},
		&model.Vote{},
		&model.MediaSeen{},
	)

	if err != nil {
		panic(err)
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
