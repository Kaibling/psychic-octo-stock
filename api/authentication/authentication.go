package authentication

import (
	"github.com/gin-gonic/gin"
)

func AddRoute(router *gin.RouterGroup) *gin.RouterGroup {
	r := router.Group("login")
	{
		r.POST("", login)
	}
	return r
}
