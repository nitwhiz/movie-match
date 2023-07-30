package model

import (
	"github.com/google/uuid"
	"time"
)

type MediaUserVotePrediction struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Index         uint64    `gorm:"autoIncrement:true" json:"index"`
	MediaID       uuid.UUID `gorm:"not null;index:idx_media_recommendation_unique,unique" json:"mediaId"`
	Media         *Media    `json:"-"`
	UserID        uuid.UUID `gorm:"not null;index:idx_media_recommendation_unique,unique" json:"userId"`
	User          *User     `json:"-"`
	PredictedVote float64   `json:"predictedVote"`
	CreatedAt     time.Time `json:"createdAt"`
}

func (r *MediaUserVotePrediction) TableName() string {
	return "media_user_vote_prediction"
}
