package users

import (
	"github.com/Kaibling/psychic-octo-stock/middleware"
	"github.com/gin-gonic/gin"
)

func AddRoute(router *gin.RouterGroup) *gin.RouterGroup {
	r := router.Group("users")
	{
		r.POST("", middleware.Authorization, userPost)
		r.GET("", middleware.Authorization, usersGet)
		r.PUT(":id", middleware.Authorization, userPut)
		r.DELETE(":id", middleware.Authorization, userDelete)
		r.GET(":id", middleware.Authorization, userGet)
		//r.GET(":userid/stocks/:stockid/:quantity", userAddStocks)
	}
	return r
}
