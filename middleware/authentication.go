package middleware

import (
	"gin/config"
	"gin/models"
	"log"
	"os"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type login struct {
	Email    string ` json:"username" binding:"required"`
	Password string ` json:"password" binding:"required"`
}

var IdentityKey = "sub"

func Authentication() *jwt.GinJWTMiddleware {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Key: []byte(os.Getenv("SECRET_KEY")),

		IdentityKey: IdentityKey,

		TokenLookup:   "header:Authorization",
		TokenHeadName: "Bearer",

		IdentityHandler: func(ctx *gin.Context) interface{} {
			var user models.User
			cliam := jwt.ExtractClaims(ctx)
			id := cliam[IdentityKey]

			db := config.GetDB()
			if err := db.First(&user, uint(id.(float64))).Error; err != nil {
				return nil
			}

			return user
		},

		Authenticator: func(ctx *gin.Context) (interface{}, error) {
			var form login
			var user models.User

			if err := ctx.ShouldBindBodyWithJSON(&form); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}

			db := config.GetDB()
			if db.Where("email = ?", form.Email).First(&user).RecordNotFound() {
				return nil, jwt.ErrFailedAuthentication
			}
			if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)).Error(); err != "" {
				return nil, jwt.ErrFailedAuthentication
			}
			return user, nil
		},

		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				claim := jwt.MapClaims{
					IdentityKey: v.ID,
				}
				return claim
			}
			return jwt.MapClaims{}
		},
		Unauthorized: func(ctx *gin.Context, code int, message string) {
			ctx.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	return authMiddleware
}
