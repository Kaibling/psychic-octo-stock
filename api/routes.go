package api

import (
	"github.com/Kaibling/psychic-octo-stock/api/stocks"
	"github.com/Kaibling/psychic-octo-stock/api/transactions"
	"github.com/Kaibling/psychic-octo-stock/api/users"
	"github.com/gin-gonic/gin"
)

func BuildRouter(router *gin.Engine) *gin.RouterGroup {
	v1 := router.Group("/v1")
	{
		users.AddRoute(v1)
		stocks.AddRoute(v1)
		transactions.AddRoute(v1)
	}
	return v1
}
