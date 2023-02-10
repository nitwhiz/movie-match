package dbutils

import (
	"github.com/nitwhiz/movie-match/server/pkg/model"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"time"
)

func GetConnection() (*gorm.DB, error) {
	dsn := "host=localhost user=root password=root dbname=movie_match sslmode=disable TimeZone=Europe/Berlin"

	verboseLogger := logger.New(
		log.New(),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Error,
			Colorful:      true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: verboseLogger,
	})

	if err != nil {
		return nil, err

	}

	return db, nil
}

func InitUsers(db *gorm.DB) error {
	users := viper.GetStringSlice("users")

	for _, userName := range users {
		log.Infof("ensuring user %s", userName)

		db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoNothing: true,
		}).Create(&model.User{
			Name: userName,
		})
	}

	return nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Media{},
		&model.Vote{},
		&model.MediaSeen{},
	)
}

func FindOrNil[ModelType interface{}](db *gorm.DB) (*ModelType, error) {
	var record ModelType

	tx := db.Limit(1).Find(&record)

	if tx.RowsAffected == 0 {
		return nil, tx.Error
	}

	return &record, tx.Error
}

func FindByIdOrNil[ModelType interface{}](db *gorm.DB, modelId string) (*ModelType, error) {
	var record ModelType

	tx := db.Where("id = ?", modelId).Limit(1).Find(&record)

	if tx.RowsAffected == 0 {
		return nil, tx.Error
	}

	return &record, tx.Error
}
