package model

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"time"
)

const MediaTypeMovie = "movie"
const MediaTitleTypeTv = "tv"

const DataSourceTMDB = "tmdb"

type Media struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ForeignID  string    `gorm:"not null;index:idx_media_unique,unique"`
	Type       string    `gorm:"not null;index:idx_media_unique,unique"`
	DataSource string    `gorm:"not null;default:'unknown'"`
	Data       datatypes.JSON
	CreatedAt  time.Time
}
