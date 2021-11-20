package api

import (
	"github.com/Kaibling/psychic-octo-stock/lib/config"
	"github.com/Kaibling/psychic-octo-stock/lib/database"
	"github.com/Kaibling/psychic-octo-stock/models"
	"github.com/Kaibling/psychic-octo-stock/repositories"
	"github.com/gin-gonic/gin"
)

func AssembleServer() *gin.Engine {
	configData := config.NewConfig()
	configData.LogEnv()
	sdb := database.NewDatabaseConnector(configData.DBUrl)
	sdb.Connect()
	sdb.Migrate(&models.User{})

	userRepo := repositories.NewUserRepository(sdb)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(injectData("userRepo", userRepo))
	BuildRouter(r)
	return r
}
func injectData(key string, data interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(key, data)
		c.Next()
	}
}
