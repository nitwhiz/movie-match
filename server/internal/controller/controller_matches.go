package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/nitwhiz/movie-match/server/internal/auth"
	"github.com/nitwhiz/movie-match/server/pkg/model"
	"gorm.io/gorm"
	"net/http"
)

type matchesParams struct {
	MediaType string `form:"type"`
}

func matchesGetAll(db *gorm.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		var matchParams matchesParams

		if err := context.ShouldBindQuery(&matchParams); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user := auth.GetJWTUser(context)

		if user == nil {
			context.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		var matches []mediaMatchResult

		voteTable := (&model.Vote{}).TableName()
		mediaTable := (&model.Media{}).TableName()

		tx := db.
			Table(voteTable + " AS a").
			Select("a.media_id AS media_id, b.user_id AS other_user_id")

		if matchParams.MediaType != "" {
			tx.Joins("JOIN "+mediaTable+" AS m ON a.media_id = m.id AND m.type = ?", matchParams.MediaType)
		}

		tx.
			Joins("JOIN "+voteTable+" AS b ON a.media_id = b.media_id AND b.user_id != ?", user.ID).
			Where("a.user_id = ? AND a.type = 'positive' AND b.type = 'positive'", user.ID)

		tx.Scan(&matches)

		context.JSON(http.StatusOK, gin.H{
			"results": matches,
		})
	}
}

func useMatches(router gin.IRouter, db *gorm.DB) {
	matchesRouter := router.Group("/matches")

	matchesRouter.GET("", matchesGetAll(db))
}
