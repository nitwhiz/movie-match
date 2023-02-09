package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/nitwhiz/movie-match/server/internal/dbutils"
	"github.com/nitwhiz/movie-match/server/internal/poster"
	"github.com/nitwhiz/movie-match/server/pkg/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
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

		media, err := dbutils.FindByIdOrNil[model.Media](db.Preload(clause.Associations), mediaParams.MediaID)

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if media == nil {
			context.Status(http.StatusNotFound)
			return
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

		media, err := dbutils.FindByIdOrNil[model.Media](db.Preload(clause.Associations), mediaPosterParams.MediaID)

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if media == nil {
			context.Status(http.StatusNotFound)
			return
		}

		fsPath, err := poster.GetPosterPath(mediaPosterParams.MediaID)

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.Status(http.StatusOK)
		context.File(fsPath)
	})
}
