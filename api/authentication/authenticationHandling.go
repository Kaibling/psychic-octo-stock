package authentication

import (
	"github.com/Kaibling/psychic-octo-stock/lib/utility"
	"github.com/Kaibling/psychic-octo-stock/models"
	"github.com/Kaibling/psychic-octo-stock/repositories"
	"github.com/gin-gonic/gin"
)

type UserLogin struct {
	Username string
	Password string
}

func login(c *gin.Context) {

	hmacSampleSecret := c.MustGet("hmacSecret").([]byte)
	var userLogin UserLogin
	c.BindJSON(&userLogin)
	userRepo := c.MustGet("userRepo").(*repositories.UserRepository)
	psHash, err := userRepo.GetPWByName(userLogin.Username)
	if err != nil {
		c.JSON(401, models.Envelope{Data: "", Message: err.Error()})
	}
	if !utility.ComparePasswords(psHash, userLogin.Password) {
		c.JSON(401, models.Envelope{Data: "", Message: "username/password incorrect"})
		c.Abort()
		return
	}

	token, erro := utility.GenerateToken(userLogin.Username, hmacSampleSecret)
	if erro != nil {
		c.JSON(500, models.Envelope{Data: "", Message: err.Error()})
	}
	c.JSON(200, models.Envelope{Data: token, Message: ""})
}
