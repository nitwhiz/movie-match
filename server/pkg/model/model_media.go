package model

import (
	"github.com/google/uuid"
	"time"
)

type MediaType = string

const MediaTypeMovie = MediaType("movie")
const MediaTypeTv = MediaType("tv")

var AvailableMediaTypes = []MediaType{
	MediaTypeMovie,
	MediaTypeTv,
}

type Media struct {
	ID          uuid.UUID   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ForeignID   string      `gorm:"index:idx_foreign_unique,unique" json:"foreignId"`
	Type        MediaType   `gorm:"type:string;index:idx_media_unique,unique;index:idx_foreign_unique,unique" json:"type"`
	Provider    string      `gorm:"default:'unknown';index:idx_foreign_unique,unique" json:"provider"`
	Title       string      `gorm:"index:idx_media_unique,unique" json:"title"`
	Summary     string      `json:"summary"`
	Genres      []Genre     `gorm:"many2many:media_genres;constraint:OnDelete:CASCADE" json:"genres"`
	Runtime     int         `gorm:"not null;default:0" json:"runtime"`
	Rating      int         `gorm:"not null;default:0" json:"rating"`
	ReleaseDate time.Time   `gorm:"type:date" json:"releaseDate"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
	Votes       []Vote      `gorm:"constraint:OnDelete:CASCADE"`
	Seen        []MediaSeen `gorm:"constraint:OnDelete:CASCADE"`
}

func (m *Media) TableName() string {
	return "media"
}
