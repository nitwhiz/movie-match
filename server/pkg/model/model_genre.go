package model

import (
	"github.com/google/uuid"
	"time"
)

type Genre struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"unique" json:"name"`
	Media     []Media   `gorm:"many2many:media_genres" json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

func (g *Genre) TableName() string {
	return "genres"
}
