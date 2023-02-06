package model

import (
	"github.com/google/uuid"
	"time"
)

type MediaSeen struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"not null;index:idx_media_seen_unique,unique;constraint:OnDelete:CASCADE" json:"userId"`
	User      User
	MediaID   uuid.UUID `gorm:"not null;index:idx_media_seen_unique,unique;constraint:OnDelete:CASCADE" json:"mediaId"`
	Media     Media
	CreatedAt time.Time `json:"createdAt"`
}

func (m *MediaSeen) TableName() string {
	return "media_seens"
}
