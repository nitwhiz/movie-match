package provider

import (
	"github.com/nitwhiz/movie-match/server/pkg/model"
	"gorm.io/gorm"
)

type MediaProvider interface {
	Init() error
	Pull(db *gorm.DB, mediaType model.MediaType, pages int) error
}

func GetMediaProviderByName(providerName string) (MediaProvider, error) {
	switch providerName {
	case tmdbProviderName:
		return NewTMDB()
	default:
		return nil, nil
	}
}
