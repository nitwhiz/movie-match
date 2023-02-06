package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nitwhiz/movie-match/server/pkg/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"path/filepath"
)

func AddMediaRetrieveAll(router gin.IRouter, db *gorm.DB) {
	router.GET("/media", func(context *gin.Context) {
		var media []model.Media

		db.Preload(clause.Associations).Limit(25).Find(&media)

		context.JSON(http.StatusOK, gin.H{
			"results": media,
		})
	})
}

func AddMediaRetrieveById(router gin.IRouter, db *gorm.DB) {
	router.GET("/media/:mediaId", func(context *gin.Context) {
		var mediaParams MediaParams

		if err := context.BindUri(&mediaParams); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var media model.Media

		if err := db.Preload(clause.Associations).Limit(1).Where("id = ?", mediaParams.MediaId).First(&media).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				context.AbortWithStatus(http.StatusNotFound)
				return
			} else {
				context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
		}

		context.JSON(http.StatusOK, media)
	})
}

func AddMediaRetrievePoster(router gin.IRouter, db *gorm.DB) {
	router.GET("/media/:mediaId/poster", func(context *gin.Context) {
		var mediaPosterParams MediaPosterParams

		if err := context.BindUri(&mediaPosterParams); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var media model.Media

		if err := db.Where("id = ?", mediaPosterParams.MediaID).Find(&media).Error; err != nil {
			context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		// todo: this is stupid
		files, err := filepath.Glob(fmt.Sprintf("/home/andy/posters/%s.*", media.ID))

		if err != nil {
			context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if len(files) == 0 {
			context.JSON(http.StatusNotFound, gin.H{"error": "0 files"})
			return
		}

		context.Status(http.StatusOK)
		context.File(files[0])
	})
}
