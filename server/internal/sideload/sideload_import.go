package sideload

import (
	"context"
	"errors"
	"fmt"
	"github.com/nitwhiz/movie-match/server/internal/config"
	"github.com/nitwhiz/movie-match/server/internal/dbutils"
	"github.com/nitwhiz/movie-match/server/internal/provider"
	"github.com/nitwhiz/movie-match/server/pkg/model"
	"github.com/ryanbradynd05/go-tmdb"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
	"os"
	"path"
	"strings"
	"sync"
)

type AdultMovieError struct {
}

func (*AdultMovieError) Error() string {
	return "adult movie"
}

func copyPoster(mediaPrefix string, originalPosterPath string, mediaSourceId string, mediaDestinationId string) (string, error) {
	ext := path.Ext(originalPosterPath)

	if ext == "" {
		return "", errors.New("unable to acquire original file name extension")
	}

	sourceFilePath := fmt.Sprintf("%s/%s-%s%s", strings.TrimRight(config.C.Sideload.Posters.FsBasePath, "/"), mediaPrefix, mediaSourceId, ext)
	posterFileName := fmt.Sprintf("%s%s", mediaDestinationId, ext)
	destinationFilePath := fmt.Sprintf("%s/%s", strings.TrimRight(config.C.Poster.FsBasePath, "/"), posterFileName)

	in, err := os.ReadFile(sourceFilePath)

	if err != nil {
		return "", err
	}

	if err := os.WriteFile(destinationFilePath, in, 0644); err != nil {
		return "", err
	}

	return posterFileName, nil
}

func processMovie(l *log.Entry, db *gorm.DB, movie *tmdb.Movie) error {
	if movie.Overview == "" {
		return errors.New("missing overview")
	}

	if movie.Adult == true {
		return &AdultMovieError{}
	}

	foreignId := fmt.Sprintf("%d", movie.ID)

	l = l.WithFields(log.Fields{
		"foreignId": foreignId,
	})

	l.Info("persisting media")

	mediaModel, err := provider.InsertMovie(db, movie, foreignId)

	if err != nil {
		return err
	}

	l = l.WithFields(log.Fields{
		"id": mediaModel.ID,
	})

	l.Info("media persisted")

	l = l.WithFields(log.Fields{
		"originalPosterPath": movie.PosterPath,
	})

	l.Info("copying poster")

	posterFileName, err := copyPoster("movie", movie.PosterPath, foreignId, mediaModel.ID.String())

	if err != nil {
		l.Warn("unable to copy poster: ", err)
	} else {
		l.Info("poster persisted")

		mediaModel.PosterFileName = posterFileName

		if err := db.Save(mediaModel).Error; err != nil {
			return err
		}
	}

	return nil
}

func processTv(l *log.Entry, db *gorm.DB, tv *tmdb.TV) error {
	if tv.Overview == "" {
		return errors.New("missing overview")
	}

	foreignId := fmt.Sprintf("%d", tv.ID)

	l = l.WithFields(log.Fields{
		"foreignId": foreignId,
	})

	l.Info("persisting media")

	mediaModel, err := provider.InsertTv(db, tv, foreignId)

	if err != nil {
		return err
	}

	l = l.WithFields(log.Fields{
		"id": mediaModel.ID,
	})

	l.Info("media persisted")

	l = l.WithFields(log.Fields{
		"originalPosterPath": tv.PosterPath,
	})

	l.Info("copying poster")

	posterFileName, err := copyPoster("tv", tv.PosterPath, foreignId, mediaModel.ID.String())

	if err != nil {
		l.Warn("unable to copy poster: ", err)
	} else {
		l.Info("poster persisted")

		mediaModel.PosterFileName = posterFileName

		if err := db.Save(mediaModel).Error; err != nil {
			return err
		}
	}

	return nil
}

func Import() error {
	if !IsConfigured() {
		return errors.New("sideload is not configured")
	}

	if !HasPostersFsAccess() {
		return errors.New("unable to access posters base path")
	}

	log.Info("connecting to mongo db")

	client, err := getMongoConnection()

	if err != nil {
		return err
	}

	log.Info("connecting to own db")

	db, err := dbutils.GetConnection()

	if err != nil {
		return err
	}

	log.Info("processing movies & tv series")

	wg := &sync.WaitGroup{}

	wg.Add(2)

	go func() {
		defer wg.Done()

		cursor, err := client.Database("tmdb").Collection("movie_details").Find(context.Background(), bson.D{})

		if err != nil {
			log.Error(err)
			return
		}

		for cursor.Next(context.Background()) {
			var obj = struct {
				ID primitive.ObjectID `bson:"_id"`
			}{}

			_ = cursor.Decode(&obj)

			l := log.WithFields(log.Fields{
				"type":    "movie",
				"mongoId": obj.ID.Hex(),
			})

			l.Info("processing media")

			var movie tmdb.Movie

			if err := cursor.Decode(&movie); err != nil {
				log.Error(err)
				return
			}

			if err := processMovie(l, db, &movie); err != nil {
				log.Warn(err)

				var adultMovieError *AdultMovieError

				if errors.As(err, &adultMovieError) {
					foundMedia, findError := dbutils.FirstModelOrNil[model.Media](db, &model.Media{ForeignID: fmt.Sprintf("%d", movie.ID)})

					if findError == nil {
						if foundMedia == nil {
							log.Info("media not present in database")
						} else {
							log.Info("removing media from database")

							tx := db.Delete(foundMedia)

							if tx.Error != nil {
								log.Error(tx.Error)
							}
						}
					} else {
						log.Error(findError)

					}
				}

				continue
			}
		}

		if err := cursor.Err(); err != nil {
			log.Error(err)
		}
	}()

	go func() {
		defer wg.Done()

		cursor, err := client.Database("tmdb").Collection("tv_series_details").Find(context.Background(), bson.D{})

		if err != nil {
			log.Error(err)
			return
		}

		for cursor.Next(context.Background()) {
			var obj = struct {
				ID primitive.ObjectID `bson:"_id"`
			}{}

			_ = cursor.Decode(&obj)

			l := log.WithFields(log.Fields{
				"type":    "tv",
				"mongoId": obj.ID.Hex(),
			})

			l.Info("processing media")

			var tv tmdb.TV

			if err := cursor.Decode(&tv); err != nil {
				log.Error(err)
				return
			}

			if err := processTv(l, db, &tv); err != nil {
				log.Warn(err)
				continue
			}
		}

		if err := cursor.Err(); err != nil {
			log.Error(err)
		}
	}()

	wg.Wait()

	return nil
}
