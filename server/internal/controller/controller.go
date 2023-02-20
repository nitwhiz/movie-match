package controller

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func useCors(router gin.IRouter) {
	corsConfig := cors.DefaultConfig()

	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Authorization")

	router.Use(cors.New(corsConfig))
}

func Init(db *gorm.DB) (*gin.Engine, error) {
	router := gin.Default()

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
