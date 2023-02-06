package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/nitwhiz/movie-match/server/pkg/model"
	"gorm.io/gorm"
	"net/http"
)

func AddUserRetrieveAll(router gin.IRouter, db *gorm.DB) {
	router.GET("/user", func(context *gin.Context) {
		var users []model.User

		db.Find(&users)

		context.JSON(http.StatusOK, gin.H{
			"results": users,
		})
	})
}
