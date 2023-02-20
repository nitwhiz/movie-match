package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/nitwhiz/movie-match/server/pkg/model"
	"gorm.io/gorm"
	"net/http"
)

func usersGetAll(db *gorm.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		var users []model.User

		db.Find(&users)

		context.JSON(http.StatusOK, gin.H{
			"results": users,
		})
	}
}

func useUsers(router gin.IRouter, db *gorm.DB) {
	userRouter := router.Group("/users")

	userRouter.GET("", usersGetAll(db))
}
