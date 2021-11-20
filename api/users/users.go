package users

import "github.com/gin-gonic/gin"

func AddRoute(router *gin.RouterGroup) *gin.RouterGroup {
	r := router.Group("users")
	{
		r.POST("", UserPost)
		r.GET("", usersGet)
		r.PUT(":id", usersPut)
		r.DELETE(":id", userDelete)
		r.GET(":id", userGet)
	}
	return r
}
