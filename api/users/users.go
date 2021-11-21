package users

import "github.com/gin-gonic/gin"

func AddRoute(router *gin.RouterGroup) *gin.RouterGroup {
	r := router.Group("users")
	{
		r.POST("", userPost)
		r.GET("", usersGet)
		r.PUT(":id", userPut)
		r.DELETE(":id", userDelete)
		r.GET(":id", userGet)
		//r.GET(":userid/stocks/:stockid/:quantity", userAddStocks)
	}
	return r
}
