package provider

import (
	"errors"
	"fmt"
	"github.com/nitwhiz/movie-match/server/pkg/model"
	"github.com/ryanbradynd05/go-tmdb"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

var tmdbProviderName = "tmdb"

type TMDBProvider struct {
	apiKey           string
	posterBaseUrl    string
	posterFsBasePath string
}

func NewTMDB() *TMDBProvider {
	return &TMDBProvider{}
}

func (p *TMDBProvider) Init() error {
	p.apiKey = viper.GetString("media_providers.tmdb.api_key")

	if p.apiKey == "" {
		return errors.New("tmdb api_key not configured")
	}

	p.posterBaseUrl = strings.TrimRight(viper.GetString("media_providers.tmdb.poster_base_url"), "/")

	if p.posterBaseUrl == "" {
		return errors.New("tmdb poster_base_url not configured")
	}

	p.posterFsBasePath = strings.TrimRight(viper.GetString("media_providers.tmdb.poster_fs_base_path"), "/")

	if p.posterFsBasePath == "" {
		return errors.New("tmdb poster_fs_base_path not configured")
	}

	return nil
}

func (p *TMDBProvider) Pull(db *gorm.DB) error {
	config := tmdb.Config{
		APIKey: p.apiKey,
	}

	api := tmdb.Init(config)

	discover, err := api.DiscoverMovie(map[string]string{
		"language":      "de-DE",
		"region":        "DE",
		"include_adult": "false",
	})

	if err != nil {
		panic(err)
	}

	for discover.Page < discover.TotalPages {
		log.Printf("processing page %d/%d ...", discover.Page, discover.TotalPages)

		for _, movieShort := range discover.Results {
			if err := (func(m tmdb.MovieShort) error {
				foreignID := fmt.Sprintf("%d", m.ID)

				log.Printf("processing movie '%s' (%s)", m.Title, foreignID)

				movie, err := api.GetMovieInfo(m.ID, map[string]string{
					"language": "de-DE",
				})

				if err != nil {
					return err
				}

				var genres []model.Genre

				for _, g := range movie.Genres {
					normalizedName := strings.TrimSpace(g.Name)

					var genre model.Genre

					if err := db.Where("name = ?", normalizedName).First(&genre).Error; err != nil {
						if errors.Is(err, gorm.ErrRecordNotFound) {
							genre.Name = g.Name
						} else {
							return err
						}
					}

					genres = append(genres, genre)
				}

				releaseDate, err := time.Parse(time.DateOnly, movie.ReleaseDate)

				mediaModel := model.Media{
					ForeignID:   foreignID,
					Type:        model.MediaTypeMovie,
					Provider:    tmdbProviderName,
					Title:       movie.Title,
					Summary:     movie.Overview,
					Genres:      genres,
					Rating:      int(math.Round(float64(movie.VoteAverage) * 10.0)),
					ReleaseDate: releaseDate,
				}

				if err := db.First(&mediaModel, "foreign_id = ?", foreignID).Error; errors.Is(err, gorm.ErrRecordNotFound) {
					if err := db.Create(&mediaModel).Error; err != nil {
						return err
					}

					log.Printf("created")
				} else {
					if err := db.Save(&mediaModel).Error; err != nil {
						return err
					}

					log.Printf("updated")
				}

				if err := os.MkdirAll(p.posterFsBasePath, 0777); err != nil {
					return err
				}

				extension := path.Ext(movie.PosterPath)

				posterFilePath := fmt.Sprintf("%s/%s%s", p.posterFsBasePath, mediaModel.ID.String(), extension)

				posterFile, err := os.Create(posterFilePath)

				if err != nil {
					return err
				}

				defer posterFile.Close()

				posterUrl := fmt.Sprintf("%s/%s", p.posterBaseUrl, strings.TrimLeft(movie.PosterPath, "/"))

				resp, err := http.Get(posterUrl)

				if err != nil {
					return err
				}

				defer resp.Body.Close()

				_, err = io.Copy(posterFile, resp.Body)

				return err
			})(movieShort); err != nil {
				return err
			}
		}

		discover, err = api.DiscoverMovie(map[string]string{
			"language": "de-DE",
			"region":   "DE",
			"page":     fmt.Sprintf("%d", discover.Page+1),
		})

		if err != nil {
			return err
		}
	}

	return nil
}
