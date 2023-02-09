package provider

import (
	"errors"
	"fmt"
	"github.com/nitwhiz/movie-match/server/internal/poster"
	"github.com/nitwhiz/movie-match/server/pkg/model"
	"github.com/ryanbradynd05/go-tmdb"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"math"
	"math/rand"
	"strings"
	"sync"
	"time"
)

var tmdbProviderName = "tmdb"

const (
	pullDiscover = iota
	pullPopular
	pullTopRated
)

const requestCooldownMin = time.Millisecond * 250
const requestCooldownMax = time.Millisecond * 800

func getRequestCooldown() time.Duration {
	return time.Duration(rand.Float64()*float64(requestCooldownMax-requestCooldownMin) + float64(requestCooldownMin))
}

type TMDBProvider struct {
	api           *tmdb.TMDb
	posterFetcher *poster.Fetcher
	language      string
	region        string
}

type genreInfo = struct {
	ID   int
	Name string
}

func NewTMDB() *TMDBProvider {
	return &TMDBProvider{}
}

func (p *TMDBProvider) Init() error {
	posterFsBasePath := strings.TrimRight(viper.GetString("media_providers.tmdb.poster_fs_base_path"), "/")

	if posterFsBasePath == "" {
		return errors.New("tmdb poster_fs_base_path not configured")
	}

	posterBaseUrl := strings.TrimRight(viper.GetString("media_providers.tmdb.poster_base_url"), "/")

	if posterBaseUrl == "" {
		return errors.New("tmdb poster_base_url not configured")
	}

	p.posterFetcher = poster.NewFetcher(posterFsBasePath, posterBaseUrl)

	p.language = viper.GetString("media_providers.tmdb.language")

	if p.language == "" {
		p.language = "en"
	}

	p.region = viper.GetString("media_providers.tmdb.region")

	if p.region == "" {
		p.region = "US"
	}

	apiKey := viper.GetString("media_providers.tmdb.api_key")

	if apiKey == "" {
		return errors.New("tmdb api_key not configured")
	}

	config := tmdb.Config{
		APIKey: apiKey,
	}

	p.api = tmdb.Init(config)

	return nil
}

func (p *TMDBProvider) Pull(db *gorm.DB) error {
	pageCount := 10

	wg := &sync.WaitGroup{}

	// tv

	go func() {
		wg.Add(1)
		defer wg.Done()

		if err := p.PullTv(db, pullDiscover, pageCount); err != nil {
			log.WithFields(log.Fields{
				"type":     model.MediaTypeTv,
				"provider": tmdbProviderName,
				"subtype":  "discover",
			}).Error(err.Error())
		}
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()

		if err := p.PullTv(db, pullPopular, pageCount); err != nil {
			log.WithFields(log.Fields{
				"type":     model.MediaTypeTv,
				"provider": tmdbProviderName,
				"subtype":  "popular",
			}).Error(err.Error())
		}
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()

		if err := p.PullTv(db, pullTopRated, pageCount); err != nil {
			log.WithFields(log.Fields{
				"type":     model.MediaTypeTv,
				"provider": tmdbProviderName,
				"subtype":  "top_rated",
			}).Error(err.Error())
		}
	}()

	// movie

	go func() {
		wg.Add(1)
		defer wg.Done()

		if err := p.PullMovie(db, pullDiscover, pageCount); err != nil {
			log.WithFields(log.Fields{
				"type":     model.MediaTypeMovie,
				"provider": tmdbProviderName,
				"subtype":  "discover",
			}).Error(err.Error())
		}
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()

		if err := p.PullMovie(db, pullPopular, pageCount); err != nil {
			log.WithFields(log.Fields{
				"type":     model.MediaTypeMovie,
				"provider": tmdbProviderName,
				"subtype":  "popular",
			}).Error(err.Error())
		}
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()

		if err := p.PullMovie(db, pullTopRated, pageCount); err != nil {
			log.WithFields(log.Fields{
				"type":     model.MediaTypeMovie,
				"provider": tmdbProviderName,
				"subtype":  "top_rated",
			}).Error(err.Error())
		}
	}()

	time.Sleep(time.Second)

	wg.Wait()

	log.Info("done")

	return nil
}

func ensureGenres(mediaGenres []genreInfo, db *gorm.DB) ([]model.Genre, error) {
	var genres []model.Genre

	for _, g := range mediaGenres {
		normalizedName := strings.TrimSpace(g.Name)

		genre := model.Genre{
			Name: normalizedName,
		}

		if err := db.FirstOrCreate(&genre, "name = ?", normalizedName).Error; err != nil {
			return genres, err
		}

		genres = append(genres, genre)
	}

	return genres, nil
}

func (p *TMDBProvider) PullMovie(db *gorm.DB, pullType int, pages int) error {
	opts := map[string]string{
		"sort_by":       "rating.desc",
		"language":      p.language,
		"region":        p.region,
		"include_adult": "false",
	}

	page := 0
	totalPages := 1

	for page < totalPages {
		opts["page"] = fmt.Sprintf("%d", page+1)

		var discover *tmdb.MoviePagedResults
		var err error

		if pullType == pullPopular {
			discover, err = p.api.GetMoviePopular(opts)
		} else if pullType == pullDiscover {
			discover, err = p.api.DiscoverMovie(opts)
		} else {
			discover, err = p.api.GetMovieTopRated(opts)
		}

		if err != nil {
			return err
		}

		page = discover.Page
		totalPages = int(math.Min(float64(pages), float64(discover.TotalPages)))

		log.WithFields(log.Fields{
			"type":       model.MediaTypeMovie,
			"provider":   tmdbProviderName,
			"page":       page,
			"totalPages": totalPages,
		}).Info("fetching discover page")

		for _, movieShort := range discover.Results {
			logFields := log.Fields{
				"type":     model.MediaTypeMovie,
				"provider": tmdbProviderName,
			}

			foreignID := fmt.Sprintf("%d", movieShort.ID)

			logFields["foreignID"] = foreignID
			logFields["title"] = movieShort.Title

			movie, err := p.api.GetMovieInfo(movieShort.ID, map[string]string{
				"language": "de",
			})

			if movie.Overview == "" {
				log.WithFields(logFields).Info("skipping: no overview")
				continue
			} else {
				log.WithFields(logFields).Info("fetching info")
			}

			if err != nil {
				return err
			}

			genres, err := ensureGenres(movie.Genres, db)

			if err != nil {
				return err
			}

			releaseDate, err := time.Parse(time.DateOnly, movie.ReleaseDate)

			if err != nil {
				return err
			}

			mediaModel := model.Media{
				ForeignID:   foreignID,
				Type:        model.MediaTypeMovie,
				Provider:    tmdbProviderName,
				Title:       movie.Title,
				Summary:     movie.Overview,
				Genres:      genres,
				Rating:      int(math.Round(float64(movie.VoteAverage) * 10.0)),
				Runtime:     int(movie.Runtime),
				ReleaseDate: releaseDate,
			}

			// todo: this tries to be atomic, but still throws errors sometimes...
			if err := db.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "provider"}, {Name: "type"}, {Name: "foreign_id"}},
				UpdateAll: true,
			}).Create(&mediaModel).Error; err != nil {
				return err
			}

			logFields["id"] = mediaModel.ID
			logFields["posterPath"] = movie.PosterPath

			if err := p.posterFetcher.Download(movie.PosterPath, mediaModel.ID.String()); err != nil {
				return err
			}

			log.WithFields(logFields).Info("done")

			time.Sleep(getRequestCooldown())
		}
	}

	return nil
}

func (p *TMDBProvider) PullTv(db *gorm.DB, pullType int, pages int) error {
	opts := map[string]string{
		"sort_by":      "vote_average.desc",
		"language":     p.language,
		"watch_region": p.region,
		"timezone":     "Europe/Berlin",
	}

	page := 0
	totalPages := 1

	for page < totalPages {
		opts["page"] = fmt.Sprintf("%d", page+1)

		var discover *tmdb.TvPagedResults
		var err error

		if pullType == pullPopular {
			discover, err = p.api.GetTvPopular(opts)
		} else if pullType == pullDiscover {
			discover, err = p.api.DiscoverTV(opts)
		} else {
			discover, err = p.api.GetTvTopRated(opts)
		}

		if err != nil {
			return err
		}

		page = discover.Page
		totalPages = int(math.Min(float64(pages), float64(discover.TotalPages)))

		log.WithFields(log.Fields{
			"type":       model.MediaTypeTv,
			"provider":   tmdbProviderName,
			"page":       page,
			"totalPages": totalPages,
		}).Info("fetching discover page")

		for _, tvShort := range discover.Results {
			logFields := log.Fields{
				"type":     model.MediaTypeTv,
				"provider": tmdbProviderName,
			}

			foreignID := fmt.Sprintf("%d", tvShort.ID)

			logFields["foreignID"] = foreignID
			logFields["title"] = tvShort.Name

			tv, err := p.api.GetTvInfo(tvShort.ID, map[string]string{
				"language": "de",
			})

			if tv.Overview == "" {
				log.WithFields(logFields).Info("skipping: no overview")
				continue
			} else {
				log.WithFields(logFields).Info("fetching info")
			}

			if err != nil {
				return err
			}

			genres, err := ensureGenres(tv.Genres, db)

			if err != nil {
				return err
			}

			releaseDate, err := time.Parse(time.DateOnly, tv.FirstAirDate)

			if err != nil {
				return err
			}

			runtime := 0

			if len(tv.EpisodeRunTime) > 0 {
				runtime = tv.EpisodeRunTime[0]
			}

			mediaModel := model.Media{
				ForeignID:   foreignID,
				Type:        model.MediaTypeTv,
				Provider:    tmdbProviderName,
				Title:       tv.Name,
				Summary:     tv.Overview,
				Genres:      genres,
				Rating:      int(math.Round(float64(tv.VoteAverage) * 10.0)),
				Runtime:     runtime,
				ReleaseDate: releaseDate,
			}

			if err := db.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "provider"}, {Name: "type"}, {Name: "foreign_id"}},
				UpdateAll: true,
			}).Create(&mediaModel).Error; err != nil {
				return err
			}

			logFields["id"] = mediaModel.ID
			logFields["posterPath"] = tv.PosterPath

			if err := p.posterFetcher.Download(tv.PosterPath, mediaModel.ID.String()); err != nil {
				return err
			}

			log.WithFields(logFields).Info("done")

			time.Sleep(getRequestCooldown())
		}
	}

	return nil
}
