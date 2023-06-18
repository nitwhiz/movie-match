package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/nitwhiz/movie-match/server/pkg/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"strings"
)

func searchMedia(db *gorm.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		q, ok := context.GetQuery("query")

		if ok {
			var medias []model.Media

			db.Preload(clause.Associations).Limit(25).Find(&medias, "POSITION(? IN LOWER(title)) > 0", strings.ToLower(q))

			context.JSON(http.StatusOK, gin.H{
				"results": medias,
			})
		}

		context.Status(http.StatusBadRequest)
	}
}

func useSearch(router gin.IRouter, db *gorm.DB) {
	searchRouter := router.Group("/search")

	searchRouter.GET("media", searchMedia(db))
}
