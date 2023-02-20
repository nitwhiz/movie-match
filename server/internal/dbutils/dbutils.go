package dbutils

import (
	"fmt"
	"github.com/nitwhiz/movie-match/server/internal/config"
	"github.com/nitwhiz/movie-match/server/pkg/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"time"
)

func GetConnection() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=movie_match sslmode=disable TimeZone=Europe/Berlin",
		config.C.Database.Host,
		config.C.Database.Port,
		config.C.Database.User,
		config.C.Database.Password,
	)

	dbLogger := logger.New(
		log.New(),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Error,
			Colorful:      true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: dbLogger,
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

func FirstOrNil[ModelType interface{}](db *gorm.DB) (*ModelType, error) {
	var record ModelType

	tx := db.Limit(1).Find(&record)

	if tx.RowsAffected == 0 {
		return nil, tx.Error
	}

	return &record, tx.Error
}

func FirstByIdOrNil[ModelType interface{}](db *gorm.DB, modelId string) (*ModelType, error) {
	var record ModelType

	tx := db.Where("id = ?", modelId).Limit(1).Find(&record)

	if tx.RowsAffected == 0 {
		return nil, tx.Error
	}

	return &record, tx.Error
}
