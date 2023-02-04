package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nitwhiz/movie-match/server/pkg/model"
	"github.com/ryanbradynd05/go-tmdb"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"net/http"
	"os"
	"time"
)

var tmdbApiKey string

func fetchDiscover(db *gorm.DB) {
	// todo: cache images

	config := tmdb.Config{
		APIKey: tmdbApiKey,
	}

	tmdbApi := tmdb.Init(config)

	discover, err := tmdbApi.DiscoverMovie(map[string]string{
		"language": "de-DE",
		"region":   "DE",
	})

	if err != nil {
		panic(err)
	}

	for discover.Page < discover.TotalPages {
		log.Printf("processing page %d/%d ...", discover.Page, discover.TotalPages)

		for _, movieShort := range discover.Results {
			(func(m tmdb.MovieShort) {
				foreignID := fmt.Sprintf("%d", m.ID)

				movie, err := tmdbApi.GetMovieInfo(m.ID, map[string]string{
					"language": "de-DE",
				})

				if err != nil {
					log.Printf("error fetching movie %s: %s", foreignID, err)
					return
				}

				movieData, err := json.Marshal(movie)

				if err != nil {
					log.Printf("error serializing movie %s: %s", foreignID, err)
				}

				var media model.Media

				if err := db.Where("foreign_id = ?", foreignID).First(&media).Error; errors.Is(err, gorm.ErrRecordNotFound) {
					movieModel := model.Media{
						ForeignID:  foreignID,
						Type:       model.MediaTypeMovie,
						DataSource: model.DataSourceTMDB,
						Data:       movieData,
					}

					if err := db.Create(&movieModel).Error; err != nil {
						log.Printf("error persisting movie %s: %s", foreignID, err)
					}
				} else {
					log.Printf("updating movie data %s", foreignID)

					media.Data = movieData

					db.Save(media)
				}
			})(movieShort)
		}

		discover, err = tmdbApi.DiscoverMovie(map[string]string{
			"language": "de-DE",
			"region":   "DE",
			"page":     fmt.Sprintf("%d", discover.Page+1),
		})

		if err != nil {
			panic(err)
		}
	}
}

type UserVoteParams struct {
	UserId   string `uri:"userId"`
	MediaId  string `json:"mediaId"`
	VoteType string `json:"voteType"`
}

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	tmdbApiKey = viper.GetString("api_keys.tmdb")

	if tmdbApiKey == "" {
		panic("missing tmdb api key")
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

	// fetchDiscover(db)

	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"hello": "world",
		})
	})

	router.GET("/media", func(context *gin.Context) {
		var media []model.Media

		db.Limit(25).Find(&media)

		context.JSON(http.StatusOK, gin.H{
			"Results": media,
		})
	})

	router.PUT("/user/:userId/vote", func(context *gin.Context) {
		var voteParams UserVoteParams

		if err := context.BindUri(&voteParams); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := context.BindJSON(&voteParams); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if voteParams.MediaId == "" {
			context.JSON(http.StatusBadRequest, gin.H{"error": "missing media id"})
			return
		}

		if voteParams.VoteType == "" {
			context.JSON(http.StatusBadRequest, gin.H{"error": "missing vote type"})
			return
		}

		var user model.User

		if err := db.Where("id = ?", voteParams.UserId).First(&user).Error; err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var media model.Media

		if err := db.Where("id = ?", voteParams.MediaId).First(&media).Error; err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		vote := model.Vote{
			User:  user,
			Media: media,
			Type:  voteParams.VoteType,
		}

		if err := db.Create(&vote).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		context.Status(http.StatusOK)
	})

	router.Run("0.0.0.0:6445")
}
