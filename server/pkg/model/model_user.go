package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"unique" json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

func (u *User) TableName() string {
	return "users"
}
