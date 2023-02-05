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
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var tmdbApiKey string
var tmdbPosterBaseUrl string

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

				mediaModel := model.Media{
					ForeignID:  foreignID,
					Type:       model.MediaTypeMovie,
					DataSource: model.DataSourceTMDB,
					Data:       movieData,
				}

				if err := db.First(&mediaModel, "foreign_id = ?", foreignID).Error; errors.Is(err, gorm.ErrRecordNotFound) {
					if err := db.Create(&mediaModel).Error; err != nil {
						log.Printf("error persisting movie %s: %s", foreignID, err)
						return
					}
				} else {
					log.Printf("updating movie data %s", foreignID)

					if err := db.Save(&mediaModel).Error; err != nil {
						log.Printf("error updating movie %s: %s", foreignID, err)
						return
					}
				}

				posterBasePath := "/home/andy/posters/"
				err = os.MkdirAll(posterBasePath, 0777)

				if err != nil {
					log.Printf("error persisting movie poster %s: %s", foreignID, err)
					return
				}

				extension := path.Ext(movie.PosterPath)

				posterFilePath := fmt.Sprintf("%s/%s%s", strings.TrimRight(posterBasePath, "/"), mediaModel.ID.String(), extension)

				posterFile, err := os.Create(posterFilePath)

				if err != nil {
					log.Printf("error persisting movie poster %s: %s", foreignID, err)
					return
				}

				defer posterFile.Close()

				posterUrl := fmt.Sprintf("%s/%s", strings.TrimRight(tmdbPosterBaseUrl, "/"), strings.TrimLeft(movie.PosterPath, "/"))

				resp, err := http.Get(posterUrl)

				if err != nil {
					log.Printf("error persisting movie poster %s: %s", foreignID, err)
					return
				}

				defer resp.Body.Close()

				_, err = io.Copy(posterFile, resp.Body)

				if err != nil {
					log.Printf("error persisting movie poster %s: %s", foreignID, err)
					return
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

type Match struct {
	MediaID     string `gorm:"media_id"`
	OtherUserID string `gorm:"other_user_id"`
}

func findMatches(db *gorm.DB, userId string) []Match {
	var matches []Match

	db.
		Table((&model.Vote{}).TableName()+" AS a").
		Select("a.media_id as media_id, b.user_id as other_user_id").
		Joins("JOIN "+(&model.Vote{}).TableName()+" AS b ON a.media_id = b.media_id AND b.user_id != ?", userId).
		Where("a.user_id = ? AND a.type = 'positive' AND b.type = 'positive'", userId).
		Scan(&matches)

	return matches
}

func isMatch(db *gorm.DB, userId string, mediaId string) bool {
	var res Match

	err := db.
		Table((&model.Vote{}).TableName()+" AS a").
		Select("a.media_id as media_id, b.user_id as other_user_id").
		Joins("JOIN "+(&model.Vote{}).TableName()+" AS b ON a.media_id = b.media_id AND b.user_id != ? AND a.media_id = ?", userId, mediaId).
		Where("a.user_id = ? AND a.type = 'positive' AND b.type = 'positive'", userId).
		First(&res).
		Error

	if err != nil {
		return false
	}

	return true
}

type UserVoteParams struct {
	UserId   string `uri:"userId"`
	MediaId  string `uri:"mediaId"`
	VoteType string `json:"voteType"`
}

type MediaPosterParams struct {
	MediaID string `uri:"mediaId"`
}

type UserSeenParams struct {
	UserId  string `uri:"userId"`
	MediaId string `uri:"mediaId"`
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

	tmdbPosterBaseUrl = viper.GetString("poster_base_urls.tmdb")

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

	router.GET("/media", func(context *gin.Context) {
		var media []model.Media

		db.Limit(25).Find(&media)

		context.JSON(http.StatusOK, gin.H{
			"Results": media,
		})
	})

	router.GET("/user", func(context *gin.Context) {
		var users []model.User

		db.Limit(25).Find(&users)

		context.JSON(http.StatusOK, gin.H{
			"Results": users,
		})
	})

	router.GET("/media/:mediaId/poster", func(context *gin.Context) {
		var mediaPosterParams MediaPosterParams

		if err := context.BindUri(&mediaPosterParams); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var media model.Media

		if err := db.Where("id = ?", mediaPosterParams.MediaID).Find(&media).Error; err != nil {
			context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		// todo: this is stupid
		files, err := filepath.Glob(fmt.Sprintf("/home/andy/posters/%s.*", media.ID))

		if err != nil {
			context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if len(files) == 0 {
			context.JSON(http.StatusNotFound, gin.H{"error": "0 files"})
			return
		}

		context.Status(http.StatusOK)
		context.File(files[0])
	})

	router.POST("/user/:userId/seen/:mediaId", func(context *gin.Context) {
		var seenParams UserSeenParams

		if err := context.BindUri(&seenParams); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user model.User

		if err := db.Where("id = ?", seenParams.UserId).First(&user).Error; err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var media model.Media

		if err := db.Where("id = ?", seenParams.MediaId).First(&media).Error; err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		seen := model.MediaSeen{
			User:  user,
			Media: media,
		}

		if err := db.Create(&seen).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		context.Status(http.StatusOK)
	})

	router.DELETE("user/:userId/seen/:mediaId", func(context *gin.Context) {
		var seenParams UserSeenParams

		if err := context.BindUri(&seenParams); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user model.User

		if err := db.Where("id = ?", seenParams.UserId).First(&user).Error; err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var media model.Media

		if err := db.Where("id = ?", seenParams.MediaId).First(&media).Error; err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var seen model.MediaSeen

		if err := db.Where("user_id = ?", user.ID).Where("media_id = ?", media.ID).First(&seen).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := db.Delete(&seen).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		context.Status(http.StatusOK)
	})

	router.PUT("/user/:userId/vote/:mediaId", func(context *gin.Context) {
		var voteParams UserVoteParams

		if err := context.BindUri(&voteParams); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := context.BindJSON(&voteParams); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		}

		if err := db.Where("user_id = ?", vote.User.ID).Where("media_id = ?", vote.Media.ID).FirstOrCreate(&vote).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		vote.Type = voteParams.VoteType

		if err := db.Save(&vote).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		context.JSON(http.StatusOK, gin.H{
			"isMatch": isMatch(db, voteParams.UserId, voteParams.MediaId),
		})
	})

	router.Run("0.0.0.0:6445")
}
