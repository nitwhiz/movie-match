package provider

import "gorm.io/gorm"

type Provider interface {
	Init() error
	Pull(db *gorm.DB) error
}
