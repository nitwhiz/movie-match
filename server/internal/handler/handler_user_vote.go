package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/nitwhiz/movie-match/server/internal/dbutils"
	"github.com/nitwhiz/movie-match/server/pkg/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

		user, err := dbutils.FindByIdOrNil[model.User](db, voteParams.UserId)

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if user == nil {
			context.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		media, err := dbutils.FindByIdOrNil[model.Media](db, voteParams.MediaId)

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if media == nil {
			context.JSON(http.StatusNotFound, gin.H{"error": "media not found"})
			return
		}

		vote := model.Vote{
			User:  user,
			Media: media,
			Type:  voteParams.VoteType,
		}

		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "media_id"}},
			UpdateAll: true,
		}).Create(&vote).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		matchTx := db.
			Table((&model.Vote{}).TableName()+" AS a").
			Select("b.user_id AS other_user_id").
			Joins("JOIN "+(&model.Vote{}).TableName()+" AS b ON a.media_id = b.media_id AND b.user_id != ? AND a.media_id = ?", voteParams.UserId, voteParams.MediaId).
			Where("a.user_id = ? AND a.type = 'positive' AND b.type = 'positive'", voteParams.UserId)

		matchResult, err := dbutils.FindOrNil[MatchResult](matchTx)

		if err != nil {
			return
		}

		context.JSON(http.StatusOK, gin.H{
			"isMatch": matchResult != nil,
		})
	})
}
