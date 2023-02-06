package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/nitwhiz/movie-match/server/pkg/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
)

func AddUserMediaRecommendedRetrieveAll(router gin.IRouter, db *gorm.DB) {
	router.GET("/user/:userId/media/recommended", func(context *gin.Context) {
		var recommendationParams MediaRecommendationParams

		if err := context.BindUri(&recommendationParams); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := context.ShouldBindQuery(&recommendationParams); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		var media []model.Media

		pageSize := 5

		db.
			Preload(clause.Associations).
			Select("m.*").
			Table((&model.Media{}).TableName()+" AS m").
			Joins("LEFT JOIN "+(&model.Vote{}).TableName()+" AS v ON m.id = v.media_id AND v.user_id = ?", recommendationParams.UserId).
			Where("v.user_id IS NULL OR v.type = 'neutral'").
			Order("m.release_date DESC, m.rating DESC").
			Offset(recommendationParams.Page * pageSize).
			Limit(pageSize).
			Find(&media)

		context.JSON(http.StatusOK, gin.H{
			"results": media,
		})
	})
}
