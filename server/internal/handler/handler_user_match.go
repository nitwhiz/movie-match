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

		var matches []MatchResult

		db.
			Table((&model.Vote{}).TableName()+" AS a").
			Select("a.media_id as media_id, b.user_id as other_user_id").
			Joins("JOIN "+(&model.Vote{}).TableName()+" AS b ON a.media_id = b.media_id AND b.user_id != ?", matchParams.UserId).
			Where("a.user_id = ? AND a.type = 'positive' AND b.type = 'positive'", matchParams.UserId).
			Scan(&matches)

		context.JSON(http.StatusOK, gin.H{
			"results": matches,
		})
	})
}
