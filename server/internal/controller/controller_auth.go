package controller

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/nitwhiz/movie-match/server/internal/auth"
	"github.com/nitwhiz/movie-match/server/pkg/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
)

func authLogout(db *gorm.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		token := jwt.GetToken(context)

		if token == "" {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		user := auth.GetJWTUser(context)

		if user == nil {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if err := db.Where(&model.UserToken{
			User:  user,
			Token: token,
		}).Delete(&model.UserToken{}).Error; err != nil {
			log.Error("Logout Error: " + err.Error())
		}

		context.Status(http.StatusNoContent)
	}
}

func authCheckTokenActivityMiddlewareFunc(db *gorm.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		_, _, _ = auth.CheckTokenActivity(context, db)
	}
}

func useAuth(router gin.IRouter, db *gorm.DB) error {
	mw, err := auth.GetJWTMiddleware(db)

	if err != nil {
		return err
	}

	// login should bypass jwt middleware
	router.POST("/auth/login", mw.LoginHandler)

	router.Use(mw.MiddlewareFunc())
	router.Use(authCheckTokenActivityMiddlewareFunc(db))

	authRouter := router.Group("/auth")

	authRouter.POST("refresh_token", mw.RefreshHandler)
	authRouter.POST("logout", authLogout(db))

	return nil
}
