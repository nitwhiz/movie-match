package model

import (
	"github.com/google/uuid"
	"time"
)

const VoteTypePositive = "positive"
const VoteTypeNegative = "negative"
const VoteTypeNeutral = "neutral"

type Vote struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"not null;index:idx_vote_unique,unique;constraint:OnDelete:CASCADE" json:"userId"`
	User      User
	MediaID   uuid.UUID `gorm:"not null;index:idx_vote_unique,unique;constraint:OnDelete:CASCADE" json:"mediaId"`
	Media     Media
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (v *Vote) TableName() string {
	return "votes"
}
