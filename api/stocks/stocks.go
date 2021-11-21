package stocks

import "github.com/gin-gonic/gin"

func AddRoute(router *gin.RouterGroup) *gin.RouterGroup {
	r := router.Group("stocks")
	{
		r.POST("/users/:userid", stockPost)
		r.GET("", stocksGet)
		r.PUT(":id", stockPut)
		r.DELETE(":id", stockDelete)
		r.GET(":id", stockGet)
	}
	return r
}
