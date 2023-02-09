package provider

import "gorm.io/gorm"

type Provider interface {
	Init() error
	Pull(db *gorm.DB) error
}

func GetByName(providerName string) Provider {
	switch providerName {
	case tmdbProviderName:
		return NewTMDB()
	default:
		return nil
	}
}
