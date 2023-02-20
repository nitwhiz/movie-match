package command

import (
	"github.com/nitwhiz/movie-match/server/internal/dbutils"
	"github.com/nitwhiz/movie-match/server/pkg/model"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func dropTable(db *gorm.DB, tblr schema.Tabler) {
	dropTableByName(db, tblr.TableName())
}

func dropTableByName(db *gorm.DB, tableName string) {
	log.WithFields(log.Fields{
		"table": tableName,
	}).Info("dropping table")

	if err := db.Exec("DROP TABLE IF EXISTS " + tableName + " CASCADE").Error; err != nil {
		log.Error(err)
	}
}

func Purge(_ *cli.Context) error {
	db, err := dbutils.GetConnection()

	if err != nil {
		return err
	}

	dropTable(db, &model.Vote{})
	dropTable(db, &model.MediaSeen{})
	dropTable(db, &model.User{})
	dropTable(db, &model.Genre{})
	dropTable(db, &model.Media{})
	dropTable(db, &model.UserToken{})

	// todo: _should_ be dropped by cascade, but isn't?
	dropTableByName(db, "media_genres")

	log.Info("done.")

	return nil
}
