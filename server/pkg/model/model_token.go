package model

import (
	"github.com/google/uuid"
	"time"
)

type UserToken struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"-"`
	UserID     uuid.UUID `gorm:"not null;index:idx_user_token_unique,unique" json:"-"`
	User       *User     `json:"-"`
	Token      string    `gorm:"not null;index:idx_user_token_unique,unique" json:"-"`
	ValidUntil time.Time `gorm:"not null" json:"-"`
}

func (t *UserToken) TableName() string {
	return "user_tokens"
}
