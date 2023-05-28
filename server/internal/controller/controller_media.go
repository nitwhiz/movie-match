package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/nitwhiz/movie-match/server/internal/auth"
	"github.com/nitwhiz/movie-match/server/internal/config"
	"github.com/nitwhiz/movie-match/server/internal/dbutils"
	"github.com/nitwhiz/movie-match/server/pkg/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"path"
)

type mediaParams struct {
	MediaID string `uri:"mediaId"`
}

type mediaPosterParams struct {
	MediaID string `uri:"mediaId"`
}

type mediaSeenParams struct {
	MediaID string `uri:"mediaId"`
}

type mediaVoteParams struct {
	MediaId  string `uri:"mediaId"`
	VoteType string `json:"voteType"`
}

type mediaMatchResult struct {
	MediaID     string `gorm:"media_id" json:"mediaId"`
	OtherUserID string `gorm:"other_user_id" json:"otherUserId"`
}

func mediaGetAll(db *gorm.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		var media []model.Media

		db.Preload(clause.Associations).Limit(25).Find(&media)

		context.JSON(http.StatusOK, gin.H{
			"results": media,
		})
	}
}

func mediaGetById(db *gorm.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		var mediaParams mediaParams

		if err := context.BindUri(&mediaParams); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		media, err := dbutils.FirstByIdOrNil[model.Media](db.Preload(clause.Associations), mediaParams.MediaID)

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if media == nil {
			context.Status(http.StatusNotFound)
			return
		}

		context.JSON(http.StatusOK, media)
	}
}

func mediaGetPosterById(db *gorm.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		var mediaPosterParams mediaPosterParams

		if err := context.BindUri(&mediaPosterParams); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		media, err := dbutils.FirstByIdOrNil[model.Media](db.Preload(clause.Associations), mediaPosterParams.MediaID)

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if media == nil {
			context.Status(http.StatusNotFound)
			return
		}

		var fsPath string

		if media.PosterFileName == "" {
			log.WithFields(log.Fields{
				"mediaId": media.ID.String(),
			}).Warn("missing poster")
		}

		fsPath = path.Join(config.C.Poster.FsBasePath, media.PosterFileName)

		context.Status(http.StatusOK)
		context.File(fsPath)
	}
}

func mediaAddSeen(db *gorm.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		var seenParams mediaSeenParams

		if err := context.BindUri(&seenParams); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user := auth.GetJWTUser(context)

		if user == nil {
			context.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		media, err := dbutils.FirstByIdOrNil[model.Media](db, seenParams.MediaID)

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if media == nil {
			context.JSON(http.StatusNotFound, gin.H{"error": "media not found"})
			return
		}

		seen := model.MediaSeen{
			User:  user,
			Media: media,
		}

		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "media_id"}},
			DoNothing: true,
		}).Create(&seen).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}

		context.Status(http.StatusOK)
	}
}

func mediaVote(db *gorm.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		var voteParams mediaVoteParams

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

		user := auth.GetJWTUser(context)

		if user == nil {
			context.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		media, err := dbutils.FirstByIdOrNil[model.Media](db, voteParams.MediaId)

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

		// todo: only update votes table
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "media_id"}},
			UpdateAll: true,
		}).Create(&vote).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		matchTx := db.
			Table((&model.Vote{}).TableName()+" AS a").
			Select("b.user_id AS other_user_id").
			Joins("JOIN "+(&model.Vote{}).TableName()+" AS b ON a.media_id = b.media_id AND b.user_id != ? AND a.media_id = ?", user.ID, voteParams.MediaId).
			Where("a.user_id = ? AND a.type = 'positive' AND b.type = 'positive'", user.ID)

		matchResult, err := dbutils.FirstOrNil[mediaMatchResult](matchTx)

		if err != nil {
			return
		}

		context.JSON(http.StatusOK, gin.H{
			"isMatch": matchResult != nil,
		})
	}
}

func useMedia(router gin.IRouter, db *gorm.DB) {
	mediaRouter := router.Group("/media")

	mediaRouter.GET("", mediaGetAll(db))
	mediaRouter.GET(":mediaId", mediaGetById(db))

	mediaRouter.GET(":mediaId/poster", mediaGetPosterById(db))

	mediaRouter.PUT(":mediaId/seen", mediaAddSeen(db))
	mediaRouter.PUT(":mediaId/vote", mediaVote(db))
}
