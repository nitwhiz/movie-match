package provider

import (
	"errors"
	"fmt"
	"github.com/nitwhiz/movie-match/server/internal/config"
	"github.com/nitwhiz/movie-match/server/pkg/model"
	"github.com/ryanbradynd05/go-tmdb"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path"
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
	pulledTitlesLock *sync.Mutex
	pulledTitles     map[string]bool
	c                *config.TMDBProviderConfig
}

type genreInfo = struct {
	ID   int
	Name string
}

func NewTMDB() (*TMDBProvider, error) {
	if config.C.MediaProviders.TMDB == nil {
		return nil, errors.New("tmdb provider is not configured")
	}

	return &TMDBProvider{
		pulledTitlesLock: &sync.Mutex{},
		pulledTitles:     map[string]bool{},
		c:                config.C.MediaProviders.TMDB,
	}, nil
}

func (p *TMDBProvider) Init() error {
	apiKey := p.c.APIKey

	if apiKey == "" {
		return errors.New("tmdb api_key not configured")
	}

	p.api = tmdb.Init(tmdb.Config{
		APIKey: apiKey,
	})

	return nil
}

func (p *TMDBProvider) downloadPoster(srcPath string, mediaId string) (string, error) {
	fsBasePath := strings.TrimRight(config.C.PosterConfig.FsBasePath, "/")

	if err := os.MkdirAll(fsBasePath, 0777); err != nil {
		return "", err
	}

	extension := path.Ext(srcPath)

	posterFilePath := fmt.Sprintf("%s/%s%s", fsBasePath, mediaId, extension)

	posterFile, err := os.Create(posterFilePath)

	if err != nil {
		return "", err
	}

	defer func(posterFile *os.File) {
		_ = posterFile.Close()
	}(posterFile)

	posterUrl := fmt.Sprintf("%s/%s", strings.TrimRight(p.c.PosterBaseUrl, "/"), strings.TrimLeft(srcPath, "/"))

	resp, err := http.Get(posterUrl)

	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", nil
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	_, err = io.Copy(posterFile, resp.Body)

	return path.Base(posterFile.Name()), err
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
		"language":      p.c.Language,
		"region":        p.c.Region,
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
				"language": p.c.Language,
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

			posterFileName, err := p.downloadPoster(movie.PosterPath, mediaModel.ID.String())

			if err != nil {
				return err
			}

			if posterFileName != "" {
				mediaModel.PosterFileName = posterFileName

				if err := db.Save(mediaModel).Error; err != nil {
					return err
				}

				logger.WithFields(log.Fields{
					"fileName": posterFileName,
				}).Info("poster persisted")
			} else {
				logger.Warn("poster not found")
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
		"language":     p.c.Language,
		"watch_region": p.c.Region,
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
				"language": p.c.Language,
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

			posterFileName, err := p.downloadPoster(tv.PosterPath, mediaModel.ID.String())

			if err != nil {
				return err
			}

			if posterFileName != "" {
				mediaModel.PosterFileName = posterFileName

				if err := db.Save(mediaModel).Error; err != nil {
					return err
				}

				logger.WithFields(log.Fields{
					"fileName": posterFileName,
				}).Info("poster persisted")
			} else {
				logger.Warn("poster not found")
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
