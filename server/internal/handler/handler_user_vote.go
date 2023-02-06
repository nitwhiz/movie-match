package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/nitwhiz/movie-match/server/pkg/model"
	"gorm.io/gorm"
	"net/http"
)

func AddUserMediaVote(router gin.IRouter, db *gorm.DB) {
	router.PUT("/user/:userId/media/:mediaId/vote", func(context *gin.Context) {
		var voteParams UserVoteParams

		if err := context.BindUri(&voteParams); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := context.BindJSON(&voteParams); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if voteParams.VoteType == "" {
			context.JSON(http.StatusBadRequest, gin.H{"error": "missing vote type"})
			return
		}

		var user model.User

		if err := db.Where("id = ?", voteParams.UserId).First(&user).Error; err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var media model.Media

		if err := db.Where("id = ?", voteParams.MediaId).First(&media).Error; err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		vote := model.Vote{
			User:  user,
			Media: media,
		}

		if err := db.Where("user_id = ?", vote.User.ID).Where("media_id = ?", vote.Media.ID).FirstOrCreate(&vote).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		vote.Type = voteParams.VoteType

		if err := db.Save(&vote).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var m MatchResult
		isMatch := false

		voteTableName := (&model.Vote{}).TableName()

		err := db.
			Table(voteTableName+" AS a").
			Select("a.media_id as media_id, b.user_id as other_user_id").
			Joins("JOIN "+voteTableName+" AS b ON a.media_id = b.media_id AND b.user_id != ? AND a.media_id = ?", voteParams.UserId, voteParams.MediaId).
			Where("a.user_id = ? AND a.type = 'positive' AND b.type = 'positive'", voteParams.UserId).
			First(&m).
			Error

		if err == nil {
			isMatch = true
		}

		context.JSON(http.StatusOK, gin.H{
			"isMatch": isMatch,
		})
	})
}
