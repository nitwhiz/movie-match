package model

import (
	"github.com/google/uuid"
	"time"
)

type MediaSeen struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"not null;index:idx_media_seen_unique,unique"`
	User      User
	MediaID   uuid.UUID `gorm:"not null;index:idx_media_seen_unique,unique"`
	Media     Media
	CreatedAt time.Time
}
