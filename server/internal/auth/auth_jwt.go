package auth

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nitwhiz/movie-match/server/internal/config"
	"github.com/nitwhiz/movie-match/server/internal/dbutils"
	"github.com/nitwhiz/movie-match/server/pkg/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var jwtIdentityKey = "userId"

func GetJWTUser(ctx *gin.Context) *model.User {
	u, exists := ctx.Get(jwtIdentityKey)

	if !exists {
		return nil
	}

	user, ok := u.(*model.User)

	if !ok {
		return nil
	}

	return user
}

func CheckTokenActivity(context *gin.Context, db *gorm.DB) (bool, *model.UserToken, error) {
	token := jwt.GetToken(context)

	if token == "" {
		context.AbortWithStatus(http.StatusUnauthorized)
		return false, nil, nil
	}

	userToken, err := dbutils.FirstOrNil[model.UserToken](db.Where(&model.UserToken{Token: token}))

	if userToken == nil || err != nil {
		context.AbortWithStatus(http.StatusUnauthorized)
		return false, nil, err
	}

	return true, userToken, nil
}

func GetJWTMiddleware(db *gorm.DB) (*jwt.GinJWTMiddleware, error) {
	// this lib doesn't provide a proper refresh token. this is alright for now but it should be replaced
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:            "movie-match",
		SigningAlgorithm: "HS512",
		Key:              []byte(config.C.Login.JWTKey),
		Timeout:          time.Hour,
		MaxRefresh:       time.Hour,
		IdentityKey:      jwtIdentityKey,
		// we still need cookies for movie/tv-show posters
		TokenLookup: "header: Authorization, cookie: jwt",
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginParams login

			if err := c.Bind(&loginParams); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			user, err := dbutils.FirstOrNil[model.User](db.Where(&model.User{Username: loginParams.Username}))

			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			if user == nil {
				return nil, jwt.ErrFailedAuthentication
			}

			if err := bcrypt.CompareHashAndPassword(
				[]byte(user.Password),
				[]byte(loginParams.Password),
			); err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			c.Set("user", user)

			return user, nil
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*model.User); ok {
				return jwt.MapClaims{
					jwtIdentityKey: v.ID,
				}
			}

			return jwt.MapClaims{}
		},
		RefreshResponse: func(c *gin.Context, code int, jwtToken string, validUntil time.Time) {
			_, userToken, _ := CheckTokenActivity(c, db)

			if userToken == nil {
				return
			}

			userToken.Token = jwtToken
			userToken.ValidUntil = validUntil

			if err := db.Save(userToken).Error; err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			c.JSON(code, gin.H{
				"token": jwtToken,
			})
		},
		LoginResponse: func(c *gin.Context, code int, jwtToken string, validUntil time.Time) {
			user, userExists := c.Get("user")

			if !userExists {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			u, ok := user.(*model.User)

			if !ok {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			userToken := model.UserToken{
				User:       u,
				Token:      jwtToken,
				ValidUntil: validUntil,
			}

			if err := db.Save(&userToken).Error; err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			c.JSON(code, gin.H{
				"token": jwtToken,
			})
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)

			user, err := dbutils.FirstOrNil[model.User](db.Where(&model.User{
				ID: uuid.MustParse(claims[jwtIdentityKey].(string)),
			}))

			if user == nil || err != nil {
				return nil
			}

			return user
		},
	})
}
