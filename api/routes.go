package api

import (
	"github.com/Kaibling/psychic-octo-stock/api/users"
	"github.com/gin-gonic/gin"
)

func BuildRouter(router *gin.Engine) *gin.RouterGroup {
	v1 := router.Group("/v1")
	{
		users.AddRoute(v1)
	}
	return v1
}
