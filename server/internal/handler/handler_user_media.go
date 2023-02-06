package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/nitwhiz/movie-match/server/pkg/model"
	"gorm.io/gorm"
	"net/http"
)

func AddUserMediaAddSeen(router gin.IRouter, db *gorm.DB) {
	router.POST("/user/:userId/seen/:mediaId", func(context *gin.Context) {
		var seenParams UserSeenParams

		if err := context.BindUri(&seenParams); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user model.User

		if err := db.Where("id = ?", seenParams.UserId).First(&user).Error; err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var media model.Media

		if err := db.Where("id = ?", seenParams.MediaId).First(&media).Error; err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		seen := model.MediaSeen{
			User:  user,
			Media: media,
		}

		if err := db.Create(&seen).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		context.Status(http.StatusOK)
	})
}

func AddUserMediaDeleteSeen(router gin.IRouter, db *gorm.DB) {
	router.DELETE("user/:userId/seen/:mediaId", func(context *gin.Context) {
		var seenParams UserSeenParams

		if err := context.BindUri(&seenParams); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user model.User

		if err := db.Where("id = ?", seenParams.UserId).First(&user).Error; err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var media model.Media

		if err := db.Where("id = ?", seenParams.MediaId).First(&media).Error; err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var seen model.MediaSeen

		if err := db.Where("user_id = ?", user.ID).Where("media_id = ?", media.ID).First(&seen).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := db.Delete(&seen).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		context.Status(http.StatusOK)
	})
}
