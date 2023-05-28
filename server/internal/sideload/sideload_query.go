package sideload

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nitwhiz/movie-match/server/pkg/model"
	"github.com/ryanbradynd05/go-tmdb"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"math"
	"os"
	"strings"
	"time"
)

type genreInfo = struct {
	ID   int
	Name string
}

func getGenres(mediaGenres []genreInfo) []model.Genre {
	genres := []model.Genre{}

	for _, g := range mediaGenres {
		genres = append(genres, model.Genre{
			Name: strings.TrimSpace(g.Name),
		})
	}

	return genres
}

func Query(movieCount int, tvCount int) error {
	if !IsConfigured() {
		return errors.New("sideload is not configured")
	}

	log.Info("connecting to mongo db")

	client, err := getMongoConnection()

	if err != nil {
		log.Error(err)
		return err
	}

	cursor, err := client.Database("tmdb").
		Collection("movie_details").
		Find(context.Background(), bson.D{})

	if err != nil {
		log.Error(err)
		return err
	}

	encoder := json.NewEncoder(os.Stdout)

	max := movieCount

	for cursor.Next(context.Background()) {
		var movie tmdb.Movie

		if err := cursor.Decode(&movie); err != nil {
			log.Error(err)
			return err
		}

		if movie.Overview == "" {
			continue
		}

		foreignId := fmt.Sprintf("%d", movie.ID)

		genres := getGenres(movie.Genres)

		if err != nil {
			log.Error(err)
			return err
		}

		releaseDate, err := time.Parse(time.DateOnly, movie.ReleaseDate)

		if err != nil {
			continue
		}

		mediaModel := model.Media{
			ForeignID:   foreignId,
			Type:        model.MediaTypeMovie,
			Provider:    "tmdb",
			Title:       movie.Title,
			Summary:     movie.Overview,
			Genres:      genres,
			Rating:      int(math.Round(float64(movie.VoteAverage) * 10.0)),
			Runtime:     int(movie.Runtime),
			ReleaseDate: releaseDate,
		}

		if err := encoder.Encode(mediaModel); err != nil {
			log.Error(err)
			return err
		}

		max--

		if max <= 0 {
			break
		}
	}

	if err := cursor.Err(); err != nil {
		log.Error(err)
	}

	cursor, err = client.Database("tmdb").
		Collection("tv_series_details").
		Find(context.Background(), bson.D{})

	if err != nil {
		log.Error(err)
		return err
	}

	max = tvCount

	for cursor.Next(context.Background()) {
		var tv tmdb.TV

		if err := cursor.Decode(&tv); err != nil {
			log.Error(err)
			return err
		}

		if tv.Overview == "" {
			continue
		}

		foreignId := fmt.Sprintf("%d", tv.ID)

		genres := getGenres(tv.Genres)

		if err != nil {
			log.Error(err)
			return err
		}

		releaseDate, err := time.Parse(time.DateOnly, tv.FirstAirDate)

		if err != nil {
			continue
		}

		runtime := 0

		if len(tv.EpisodeRunTime) > 0 {
			runtime = tv.EpisodeRunTime[0]
		}

		mediaModel := model.Media{
			ForeignID:   foreignId,
			Type:        model.MediaTypeTv,
			Provider:    "tmdb",
			Title:       tv.Name,
			Summary:     tv.Overview,
			Genres:      genres,
			Rating:      int(math.Round(float64(tv.VoteAverage) * 10.0)),
			Runtime:     runtime,
			ReleaseDate: releaseDate,
		}

		if err := encoder.Encode(mediaModel); err != nil {
			log.Error(err)
			return err
		}

		max--

		if max <= 0 {
			break
		}
	}

	if err := cursor.Err(); err != nil {
		log.Error(err)
	}

	return nil
}
