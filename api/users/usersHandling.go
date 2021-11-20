package users

import (
	"github.com/Kaibling/psychic-octo-stock/lib/utility"
	"github.com/Kaibling/psychic-octo-stock/models"
	"github.com/Kaibling/psychic-octo-stock/repositories"
	"github.com/gin-gonic/gin"
)

func UserPost(c *gin.Context) {
	var newUser models.User
	c.BindJSON(&newUser)
	newUser.Password = utility.HashPassword(newUser.Password)
	userRepo := c.MustGet("userRepo").(*repositories.UserRepository)
	if err := userRepo.AddUser(&newUser); err != nil {
		c.JSON(500, models.Envelope{Data: "", Message: err.Error()})
		return
	}
	//todo proper return schema
	newUser.Password = ""
	env := models.Envelope{Data: newUser, Message: ""}

	c.JSON(201, env)
}
func usersGet(c *gin.Context)   {}
func usersPut(c *gin.Context)   {}
func userDelete(c *gin.Context) {}
func userGet(c *gin.Context)    {}