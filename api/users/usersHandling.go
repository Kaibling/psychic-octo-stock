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
		c.JSON(err.HttpStatus(), models.Envelope{Data: "", Message: err.Error()})
		return
	}
	//todo proper return schema
	newUser.Password = ""
	env := models.Envelope{Data: newUser, Message: ""}
	c.JSON(201, env)
}
func usersGet(c *gin.Context) {
	userRepo := c.MustGet("userRepo").(*repositories.UserRepository)
	userList, err := userRepo.GetAllUser()
	if err != nil {
		c.JSON(err.HttpStatus(), models.Envelope{Data: "", Message: err.Error()})
		return
	}
	env := models.Envelope{Data: userList, Message: ""}
	c.JSON(200, env)
}
func usersPut(c *gin.Context) {
	userID := c.Param("id")
	var updateUser map[string]interface{}
	c.BindJSON(&updateUser)

	updateUser["ID"] = userID
	userRepo := c.MustGet("userRepo").(*repositories.UserRepository)

	userRepo.UpdateWithMap(updateUser)
	loadedUser, err := userRepo.GetUserByID(userID)
	if err != nil {

		c.JSON(err.HttpStatus(), models.Envelope{Data: "", Message: err.Error()})
		return
	}
	env := models.Envelope{Data: loadedUser, Message: ""}
	c.JSON(200, env)
}
func userDelete(c *gin.Context) {
	userID := c.Param("id")
	userRepo := c.MustGet("userRepo").(*repositories.UserRepository)
	loadedUser, err := userRepo.GetUserByID(userID)
	if err != nil {
		c.JSON(err.HttpStatus(), models.Envelope{Data: "", Message: err.Error()})
		return
	}
	if err := userRepo.DeleteUserByID(loadedUser); err != nil {
		c.JSON(err.HttpStatus(), models.Envelope{Data: "", Message: err.Error()})
		return
	}
	c.JSON(204, nil)
}
func userGet(c *gin.Context) {}
