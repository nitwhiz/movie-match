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

type pullType string

const pullTypeDiscover = pullType("discover")
const pullTypePopular = pullType("popular")
const pullTypeTopRated = pullType("top_rated")

const requestCooldownMin = time.Millisecond * 250
const requestCooldownMax = time.Millisecond * 800

func getRequestCooldown() time.Duration {
	return time.Duration(rand.Float64()*float64(requestCooldownMax-requestCooldownMin) + float64(requestCooldownMin))
}

type TMDBProvider struct {
	api              *tmdb.TMDb
	posterFetcher    *poster.Fetcher
	language         string
	region           string
	pulledTitlesLock *sync.Mutex
	pulledTitles     map[string]bool
}

type genreInfo = struct {
	ID   int
	Name string
}

func NewTMDB() *TMDBProvider {
	return &TMDBProvider{
		pulledTitlesLock: &sync.Mutex{},
		pulledTitles:     map[string]bool{},
	}
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

func (p *TMDBProvider) StartPull(wg *sync.WaitGroup, db *gorm.DB, mediaType model.MediaType, pullType pullType, pages int) {
	defer wg.Done()

	l := log.WithFields(log.Fields{
		"type":     mediaType,
		"provider": tmdbProviderName,
		"subtype":  pullType,
	})

	if mediaType == model.MediaTypeTv {
		if err := p.PullTv(l, db, pullType, pages); err != nil {
			l.Error(err.Error())
		}
	} else if mediaType == model.MediaTypeMovie {
		if err := p.PullMovie(l, db, pullType, pages); err != nil {
			l.Error(err.Error())
		}
	}
}

func (p *TMDBProvider) Pull(db *gorm.DB, mediaType model.MediaType, pages int) error {
	wg := &sync.WaitGroup{}

	wg.Add(3)

	go p.StartPull(wg, db, mediaType, pullTypeDiscover, pages)
	go p.StartPull(wg, db, mediaType, pullTypePopular, pages)
	go p.StartPull(wg, db, mediaType, pullTypeTopRated, pages)

	wg.Wait()

	log.WithFields(log.Fields{
		"type": mediaType,
	}).Info("all media pulled")

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

func (p *TMDBProvider) PullMovie(logger *log.Entry, db *gorm.DB, pullType pullType, pages int) error {
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

		if pullType == pullTypePopular {
			discover, err = p.api.GetMoviePopular(opts)
		} else if pullType == pullTypeDiscover {
			discover, err = p.api.DiscoverMovie(opts)
		} else {
			discover, err = p.api.GetMovieTopRated(opts)
		}

		if err != nil {
			return err
		}

		page = discover.Page
		totalPages = int(math.Min(float64(pages), float64(discover.TotalPages)))

		logger = logger.WithFields(log.Fields{
			"page":       page,
			"totalPages": totalPages,
		})

		logger.Info("processing page")

		for _, movieShort := range discover.Results {
			foreignID := fmt.Sprintf("%d", movieShort.ID)

			logger := logger.WithFields(log.Fields{
				"foreignID": foreignID,
				"title":     movieShort.Title,
			})

			p.pulledTitlesLock.Lock()

			if p.pulledTitles[foreignID] {
				p.pulledTitlesLock.Unlock()
				logger.Info("skipping: title already pulled")
				continue
			}

			p.pulledTitles[foreignID] = true
			p.pulledTitlesLock.Unlock()

			if movieShort.Overview == "" {
				logger.Info("skipping: no short overview")
				continue
			}

			movie, err := p.api.GetMovieInfo(movieShort.ID, map[string]string{
				"language": "de",
			})

			if err != nil {
				logger.Error(err)
				return err
			}

			if movie.Overview == "" {
				logger.Info("skipping: no overview")
				continue
			}

			logger.Info("processing media")

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

			logger = logger.WithFields(log.Fields{
				"id": mediaModel.ID,
			})

			logger.Info("downloading poster")

			if err := p.posterFetcher.Download(movie.PosterPath, mediaModel.ID.String()); err != nil {
				return err
			}

			requestCooldown := getRequestCooldown()

			logger.WithFields(log.Fields{
				"requestCooldown": fmt.Sprintf("%dms", int(requestCooldown/time.Millisecond)),
			}).Info("media persisted")

			time.Sleep(requestCooldown)
		}
	}

	return nil
}

func (p *TMDBProvider) PullTv(logger *log.Entry, db *gorm.DB, pullType pullType, pages int) error {
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

		if pullType == pullTypePopular {
			discover, err = p.api.GetTvPopular(opts)
		} else if pullType == pullTypeDiscover {
			discover, err = p.api.DiscoverTV(opts)
		} else {
			discover, err = p.api.GetTvTopRated(opts)
		}

		if err != nil {
			return err
		}

		page = discover.Page
		totalPages = int(math.Min(float64(pages), float64(discover.TotalPages)))

		logger = logger.WithFields(log.Fields{
			"page":       page,
			"totalPages": totalPages,
		})

		logger.Info("processing page")

		for _, tvShort := range discover.Results {
			foreignID := fmt.Sprintf("%d", tvShort.ID)

			logger := logger.WithFields(log.Fields{
				"foreignID": foreignID,
				"title":     tvShort.Name,
			})

			p.pulledTitlesLock.Lock()

			if p.pulledTitles[foreignID] {
				p.pulledTitlesLock.Unlock()
				logger.Info("skipping: title already pulled")
				continue
			}

			p.pulledTitles[foreignID] = true
			p.pulledTitlesLock.Unlock()

			if tvShort.Overview == "" {
				logger.Info("skipping: no short overview")
				continue
			}

			tv, err := p.api.GetTvInfo(tvShort.ID, map[string]string{
				"language": "de",
			})

			if err != nil {
				logger.Error(err)
				return err
			}

			if tv.Overview == "" {
				logger.Info("skipping: no overview")
				continue
			}

			logger.Info("processing media")

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

			logger = logger.WithFields(log.Fields{
				"id": mediaModel.ID,
			})

			logger.Info("downloading poster")

			if err := p.posterFetcher.Download(tv.PosterPath, mediaModel.ID.String()); err != nil {
				return err
			}

			requestCooldown := getRequestCooldown()

			logger.WithFields(log.Fields{
				"requestCooldown": fmt.Sprintf("%dms", int(requestCooldown/time.Millisecond)),
			}).Info("media persisted")

			time.Sleep(requestCooldown)
		}
	}

	return nil
}
