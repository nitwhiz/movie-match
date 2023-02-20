package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/nitwhiz/movie-match/server/internal/auth"
	"github.com/nitwhiz/movie-match/server/pkg/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
)

type meRecommendationParams struct {
	Page int `form:"page"`
}

func meGetAllRecommendations(db *gorm.DB) gin.HandlerFunc {
	// todo: we should go jsonapi asap

	return func(context *gin.Context) {
		var recommendationParams meRecommendationParams

		if err := context.BindUri(&recommendationParams); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := context.ShouldBindQuery(&recommendationParams); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		user := auth.GetJWTUser(context)

		if user == nil {
			context.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		var media []model.Media

		pageSize := 10

		db.
			Preload(clause.Associations).
			Select("m.*").
			Table((&model.Media{}).TableName()+" AS m").
			Joins("LEFT JOIN "+(&model.Vote{}).TableName()+" AS v ON m.id = v.media_id AND v.user_id = ?", user.ID).
			Joins("LEFT JOIN "+(&model.MediaSeen{}).TableName()+" AS s ON m.id = s.media_id AND s.user_id = ?", user.ID).
			Where("s.id IS NULL").
			Where("v.id IS NULL OR v.type = 'neutral'").
			Order("m.release_date DESC, m.rating DESC").
			Offset(recommendationParams.Page * pageSize).
			Limit(pageSize).
			Find(&media)

		context.JSON(http.StatusOK, gin.H{
			"results": media,
		})
	}
}

func meGet() gin.HandlerFunc {
	return func(context *gin.Context) {
		user := auth.GetJWTUser(context)

		if user == nil {
			context.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		context.JSON(http.StatusOK, gin.H{
			"user": user,
		})
	}
}

func useMe(router gin.IRouter, db *gorm.DB) {
	meRouter := router.Group("/me")

	meRouter.GET("", meGet())
	meRouter.GET("recommended", meGetAllRecommendations(db))
}
