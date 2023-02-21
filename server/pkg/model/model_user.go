package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Username    string    `gorm:"index:idx_username_unique,unique;index:idx_unique_name,unique" json:"username"`
	Password    string    `json:"-"`
	DisplayName string    `gorm:"index:idx_unique_name,unique" json:"displayName"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (u *User) TableName() string {
	return "users"
}
