package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nitwhiz/movie-match/server/internal/auth"
	"github.com/nitwhiz/movie-match/server/internal/dbutils"
	"github.com/nitwhiz/movie-match/server/pkg/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
)

type meRecommendationParams struct {
	BelowScore string `form:"belowScore"`
}

type userMedia struct {
	model.Media
	Score    string `json:"score"`
	Seen     bool   `json:"seen"`
	VoteType string `gorm:"vote_type" json:"voteType"`
}

type voteByMediaIdParams struct {
	MediaID string `uri:"mediaId"`
}

func meGetRecommendations(db *gorm.DB) gin.HandlerFunc {
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

		var media []*userMedia

		pageSize := 25

		db.
			Select("m.*, COALESCE(v2.type, 'none') AS vote_type, c.score").
			Table((&model.Media{}).TableName()+" AS m").
			Joins("LEFT JOIN "+(&model.Vote{}).TableName()+" v ON m.id = v.media_id AND v.user_id != ?", user.ID).
			Joins("LEFT JOIN "+(&model.Vote{}).TableName()+" v2 ON m.id = v2.media_id AND v2.user_id = ?", user.ID).
			Joins(`CROSS JOIN LATERAL (VALUES (
				LEAST(
						   (
							   ABS(HASHTEXT(m.id::text)) / 2147483647.0
							   )
						   * (
							   CASE
								   WHEN v.type = 'positive' THEN 33
								   WHEN v.type = 'neutral' THEN 60
								   WHEN v.type = 'negative' THEN 75
								   WHEN v.type IS NULL THEN 100
								   END
							   )
					   + (
							   CASE
								   WHEN v.type = 'positive' THEN 67
								   WHEN v.type = 'neutral' THEN 40
								   WHEN v.type = 'negative' THEN 0
								   WHEN v.type IS NULL THEN 0
								   END
							)
				, 100.0
				)
			)) c(score)`).
			Where("(v2.type IS NULL) AND c.score < ?", recommendationParams.BelowScore).
			Order("c.score DESC, v.type DESC").
			Limit(pageSize).
			Find(&media)

		// begin workaround for incorrect preloading
		// see https://github.com/go-gorm/gorm/pull/6067

		// todo: all of this can be thrown out by using jsonapi relationships
		// -> let client request what it needs to display this information

		uniqueMediaIds := map[uuid.UUID]struct{}{}
		mediaMap := map[uuid.UUID]*userMedia{}

		for _, m := range media {
			func(m *userMedia) {
				uniqueMediaIds[m.ID] = struct{}{}

				mediaMap[m.ID] = m
				mediaMap[m.ID].Genres = []model.Genre{}

				seenMedia, err := dbutils.FirstModelOrNil[model.MediaSeen](db, &model.MediaSeen{
					MediaID: m.ID,
					UserID:  user.ID,
				})

				if err != nil {
					log.Error(err)
					return
				}

				if seenMedia != nil {
					mediaMap[m.ID].Seen = true
				}
			}(m)
		}

		mediaIds := make([]uuid.UUID, len(uniqueMediaIds))

		i := 0

		for mid := range uniqueMediaIds {
			mediaIds[i] = mid
			i += 1
		}

		var mediaGenres []struct {
			MediaID uuid.UUID `gorm:"media_id"`
			GenreID uuid.UUID `gorm:"genre_id"`
		}

		db.Table("media_genres").Where("media_id IN ?", mediaIds).Find(&mediaGenres)

		uniqueGenreIds := map[uuid.UUID]struct{}{}

		for _, g := range mediaGenres {
			uniqueGenreIds[g.GenreID] = struct{}{}
		}

		genreIds := make([]uuid.UUID, len(uniqueGenreIds))

		i = 0

		for gid := range uniqueGenreIds {
			genreIds[i] = gid
			i += 1
		}

		var genres []*model.Genre

		db.Where("id IN ?", genreIds).Find(&genres)

		genreMap := map[uuid.UUID]*model.Genre{}

		for _, g := range genres {
			genreMap[g.ID] = g
		}

		for _, m := range mediaGenres {
			mediaMap[m.MediaID].Genres = append(mediaMap[m.MediaID].Genres, *(genreMap[m.GenreID]))
		}

		// end workaround for incorrect preloading

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

func meGetVoteByMediaId(db *gorm.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		user := auth.GetJWTUser(context)

		if user == nil {
			context.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		var voteParams voteByMediaIdParams

		if err := context.BindUri(&voteParams); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		media, err := dbutils.FirstByIdOrNil[model.Media](db.Preload(clause.Associations), voteParams.MediaID)

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if media == nil {
			context.Status(http.StatusNotFound)
			return
		}

		voteType := model.VoteTypeNone

		for _, v := range media.Votes {
			if v.UserID.String() == user.ID.String() {
				voteType = v.Type
				break
			}
		}

		context.JSON(http.StatusOK, gin.H{
			"voteType": voteType,
		})
	}
}

func meGetVotes(db *gorm.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		user := auth.GetJWTUser(context)

		if user == nil {
			context.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		votes := []model.Vote{}

		db.Where("user_id = ?", user.ID).Where("type != ?", model.VoteTypeNeutral).Order("updated_at DESC").Find(&votes)

		context.JSON(http.StatusOK, gin.H{
			"results": votes,
		})
	}
}

func useMe(router gin.IRouter, db *gorm.DB) {
	meRouter := router.Group("/me")

	meRouter.GET("", meGet())
	meRouter.GET("recommended", meGetRecommendations(db))
	meRouter.GET("votes", meGetVotes(db))
	meRouter.GET("vote/:mediaId", meGetVoteByMediaId(db))
}
