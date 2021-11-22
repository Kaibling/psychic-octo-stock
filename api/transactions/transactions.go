package transactions

import (
	"github.com/Kaibling/psychic-octo-stock/middleware"
	"github.com/gin-gonic/gin"
)

func AddRoute(router *gin.RouterGroup) *gin.RouterGroup {
	r := router.Group("transactions")
	{
		r.POST("", middleware.Authorization, transactionPost)
		r.GET("", middleware.Authorization, transactionsGet)
		r.DELETE(":id", middleware.Authorization, transactionDelete)
		r.GET(":id", middleware.Authorization, transactionGet)
	}
	return r
}
