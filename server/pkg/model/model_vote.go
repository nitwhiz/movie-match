package model

import (
	"github.com/google/uuid"
	"time"
)

const VoteTypePositive = "positive"
const VoteTypeNegative = "negative"
const VoteTypeNeutral = "neutral"

type Vote struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"not null;index:idx_vote_unique,unique"`
	User      User
	MediaID   uuid.UUID `gorm:"not null;index:idx_vote_unique,unique"`
	Media     Media
	Type      string `gorm:"not null"`
	CreatedAt time.Time
}
