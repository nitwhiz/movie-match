package controller

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nitwhiz/movie-match/server/internal/config"
	"gorm.io/gorm"
)

func useCors(router gin.IRouter) {
	corsConfig := cors.DefaultConfig()

	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Authorization")

	router.Use(cors.New(corsConfig))
}

func useLoggerAndRecovery(router gin.IRouter) {
	router.Use(loggerMiddleware(), gin.Recovery())
}

func Init(db *gorm.DB) (*gin.Engine, error) {
	ginMode := gin.ReleaseMode

	if config.C.Debug {
		ginMode = gin.DebugMode
	}

	gin.SetMode(ginMode)

	router := gin.New()

	useLoggerAndRecovery(router)
	useCors(router)

	if err := useAuth(router, db); err != nil {
		return nil, err
	}

	useMedia(router, db)
	useUsers(router, db)
	useMatches(router, db)
	useMe(router, db)

	return router, nil
}
