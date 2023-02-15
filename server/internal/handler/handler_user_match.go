package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/nitwhiz/movie-match/server/pkg/model"
	"gorm.io/gorm"
	"net/http"
)

func AddUserMatchRetrieveAll(router gin.IRouter, db *gorm.DB) {
	router.GET("/user/:userId/match", func(context *gin.Context) {
		var matchParams UserMatchParams

		if err := context.BindUri(&matchParams); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := context.ShouldBindQuery(&matchParams); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var matches []MatchResult

		voteTable := (&model.Vote{}).TableName()
		mediaTable := (&model.Media{}).TableName()

		tx := db.
			Table(voteTable + " AS a").
			Select("a.media_id AS media_id, b.user_id AS other_user_id")

		if matchParams.MediaType != "" {
			tx.Joins("JOIN "+mediaTable+" AS m ON a.media_id = m.id AND m.type = ?", matchParams.MediaType)
		}

		tx.
			Joins("JOIN "+voteTable+" AS b ON a.media_id = b.media_id AND b.user_id != ?", matchParams.UserId).
			Where("a.user_id = ? AND a.type = 'positive' AND b.type = 'positive'", matchParams.UserId)

		tx.Scan(&matches)

		context.JSON(http.StatusOK, gin.H{
			"results": matches,
		})
	})
}
