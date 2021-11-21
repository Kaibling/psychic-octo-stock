package transactions

import "github.com/gin-gonic/gin"

func AddRoute(router *gin.RouterGroup) *gin.RouterGroup {
	r := router.Group("transactions")
	{
		r.POST("", transactionPost)
		r.GET("", transactionsGet)
		r.DELETE(":id", transactionDelete)
		r.GET(":id", transactionGet)
	}
	return r
}
