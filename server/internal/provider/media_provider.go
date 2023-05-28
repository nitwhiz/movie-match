package provider

import (
	"github.com/nitwhiz/movie-match/server/pkg/model"
	"gorm.io/gorm"
	"strings"
)

type MediaProvider interface {
	Init() error
	Pull(db *gorm.DB, mediaType model.MediaType, pages int) error
}

func EnsureGenres(mediaGenres []genreInfo, db *gorm.DB) ([]model.Genre, error) {
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

func GetMediaProviderByName(providerName string) (MediaProvider, error) {
	switch providerName {
	case tmdbProviderName:
		return NewTMDB()
	default:
		return nil, nil
	}
}
