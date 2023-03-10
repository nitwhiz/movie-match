package dbutils

import (
	"fmt"
	"github.com/nitwhiz/movie-match/server/internal/config"
	"github.com/nitwhiz/movie-match/server/pkg/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetConnection() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=movie_match sslmode=disable TimeZone=Europe/Berlin",
		config.C.Database.Host,
		config.C.Database.Port,
		config.C.Database.User,
		config.C.Database.Password,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: NewLogger(config.C.Debug),
	})

	if err != nil {
		return nil, err

	}

	return db, nil
}

func InitUsers(db *gorm.DB) error {
	var usernames []string

	for _, user := range config.C.Login.Users {
		log.Infof("ensuring user %s", user.Username)

		db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "username"}, {Name: "display_name"}},
			UpdateAll: true,
		}).Create(&model.User{
			DisplayName: user.DisplayName,
			Username:    user.Username,
			Password:    user.Password,
		})

		usernames = append(usernames, user.Username)
	}

	if len(usernames) > 0 {
		if err := db.Not(map[string]interface{}{"username": usernames}).Delete(&model.User{}).Error; err != nil {
			return err
		}
	}

	return nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.UserToken{},
		&model.Media{},
		&model.Vote{},
		&model.MediaSeen{},
	)
}

func first[ModelType interface{}](db *gorm.DB) (*ModelType, error) {
	var record ModelType

	tx := db.Limit(1).Find(&record)

	if tx.RowsAffected == 0 {
		return nil, tx.Error
	}

	return &record, tx.Error
}

func FirstOrNil[ModelType interface{}](db *gorm.DB) (*ModelType, error) {
	return first[ModelType](db)
}

func FirstModelOrNil[ModelType interface{}](db *gorm.DB, where *ModelType) (*ModelType, error) {
	return FirstOrNil[ModelType](db.Where(where))
}

func FirstByIdOrNil[ModelType interface{}](db *gorm.DB, modelId string) (*ModelType, error) {
	return first[ModelType](db.Where("id = ?", modelId))
}
