package model

import (
	"github.com/google/uuid"
	"time"
)

const MediaTypeMovie = "movie"
const MediaTypeTv = "tv"

type Media struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ForeignID   string    `gorm:"index:idx_media_unique,unique" json:"foreignId"`
	Type        string    `gorm:"index:idx_media_unique,unique" json:"type"`
	Provider    string    `gorm:"default:'unknown';index:idx_media_unique,unique" json:"provider"`
	Title       string    `json:"title"`
	Summary     string    `json:"summary"`
	Genres      []Genre   `gorm:"many2many:media_genres" json:"genres"`
	Rating      int       `gorm:"not null;default:0" json:"rating"`
	ReleaseDate time.Time `json:"releaseDate"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (m *Media) TableName() string {
	return "media"
}
