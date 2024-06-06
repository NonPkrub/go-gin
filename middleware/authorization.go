package middleware

import (
	"gin/models"
	"net/http"

	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
)

func Authorize() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, ok := ctx.Get("sub")
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err ": "Unauthorized"})
			return
		}

		enforcer := casbin.NewEnforcer("config/acl_model.conf", "config/policy.csv")
		ok = enforcer.Enforce(user.(*models.User).Role, ctx.Request.URL.Path, ctx.Request.Method)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"err ": "you are not allowed to access this resource"})
			return
		}
	}
}
