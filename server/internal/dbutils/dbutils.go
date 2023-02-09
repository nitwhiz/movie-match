package dbutils

import (
	"gorm.io/gorm"
)

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
