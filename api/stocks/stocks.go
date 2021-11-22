package stocks

import (
	"github.com/Kaibling/psychic-octo-stock/middleware"
	"github.com/gin-gonic/gin"
)

func AddRoute(router *gin.RouterGroup) *gin.RouterGroup {
	r := router.Group("stocks")
	{
		r.POST("/users/:userid", middleware.Authorization, stockPost)
		r.GET("", middleware.Authorization, stocksGet)
		r.PUT(":id", middleware.Authorization, stockPut)
		r.DELETE(":id", middleware.Authorization, stockDelete)
		r.GET(":id", middleware.Authorization, stockGet)
	}
	return r
}
