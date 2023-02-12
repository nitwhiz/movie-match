package provider

import (
	"github.com/nitwhiz/movie-match/server/pkg/model"
	"gorm.io/gorm"
)

type Provider interface {
	Init() error
	Pull(db *gorm.DB, mediaType model.MediaType, pages int) error
}

func GetByName(providerName string) Provider {
	switch providerName {
	case tmdbProviderName:
		return NewTMDB()
	default:
		return nil
	}
}
