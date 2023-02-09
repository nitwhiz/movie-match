package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/nitwhiz/movie-match/server/internal/dbutils"
	"github.com/nitwhiz/movie-match/server/pkg/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
)

func AddUserMediaAddSeen(router gin.IRouter, db *gorm.DB) {
	router.PUT("/user/:userId/media/:mediaId/seen", func(context *gin.Context) {
		var seenParams UserSeenParams

		if err := context.BindUri(&seenParams); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := dbutils.FindByIdOrNil[model.User](db, seenParams.UserID)

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if user == nil {
			context.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		media, err := dbutils.FindByIdOrNil[model.Media](db, seenParams.MediaID)

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if media == nil {
			context.JSON(http.StatusNotFound, gin.H{"error": "media not found"})
			return
		}

		seen := model.MediaSeen{
			User:  user,
			Media: media,
		}

		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "media_id"}},
			DoNothing: true,
		}).Create(&seen).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}

		context.Status(http.StatusOK)
	})
}
