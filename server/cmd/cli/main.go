package main

import (
	"github.com/nitwhiz/movie-match/server/internal/config"
	"github.com/nitwhiz/movie-match/server/internal/provider"
	"github.com/nitwhiz/movie-match/server/pkg/model"
	"github.com/spf13/viper"
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

	users := viper.GetStringSlice("users")

	for _, userName := range users {
		(func(name string) {
			log.Printf("creating user %s", name)

			u := model.User{
				Name: name,
			}

			db.Where("name = ?", u.Name).FirstOrCreate(&u)
		})(userName)
	}

	p := provider.NewTMDB()

	if err := p.Init(); err != nil {
		panic(err)
	}

	if err := p.Pull(db); err != nil {
		panic(err)
	}
}
