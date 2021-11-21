package stocks

import (
	"github.com/Kaibling/psychic-octo-stock/lib/utility"
	"github.com/Kaibling/psychic-octo-stock/models"
	"github.com/Kaibling/psychic-octo-stock/repositories"
	"github.com/gin-gonic/gin"
)

func stockPost(c *gin.Context) {
	userID, err := utility.GetParam("userid", c)
	if err != nil {
		c.JSON(err.HttpStatus(), models.Envelope{Data: "", Message: err.Error()})
		return
	}
	userRepo := c.MustGet("userRepo").(*repositories.UserRepository)
	_, err = userRepo.GetByID(userID)
	if err != nil {
		c.JSON(err.HttpStatus(), models.Envelope{Data: "", Message: err.Error()})
		return
	}
	var newStock models.Stock
	c.BindJSON(&newStock)
	stockRepo := c.MustGet("stockRepo").(*repositories.StockRepository)
	if err := stockRepo.Add(&newStock); err != nil {
		c.JSON(err.HttpStatus(), models.Envelope{Data: "", Message: err.Error()})
		return
	}
	//add ownership to an user
	stockRepo.AddStockToUser(newStock.ID, userID, newStock.Quantity)

	env := models.Envelope{Data: newStock, Message: ""}
	c.JSON(201, env)
}
func stocksGet(c *gin.Context) {
	stockRepo := c.MustGet("stockRepo").(*repositories.StockRepository)
	userList, err := stockRepo.GetAll()
	if err != nil {
		c.JSON(err.HttpStatus(), models.Envelope{Data: "", Message: err.Error()})
		return
	}
	env := models.Envelope{Data: userList, Message: ""}
	c.JSON(200, env)
}
func stockPut(c *gin.Context) {
	stockID := c.Param("id")
	var updateStock map[string]interface{}
	c.BindJSON(&updateStock)

	updateStock["ID"] = stockID
	stockRepo := c.MustGet("stockRepo").(*repositories.StockRepository)

	stockRepo.UpdateWithMap(updateStock)
	loadedUser, err := stockRepo.GetByID(stockID)
	if err != nil {

		c.JSON(err.HttpStatus(), models.Envelope{Data: "", Message: err.Error()})
		return
	}
	env := models.Envelope{Data: loadedUser, Message: ""}
	c.JSON(200, env)
}
func stockDelete(c *gin.Context) {
	stockID := c.Param("id")
	stockRepo := c.MustGet("stockRepo").(*repositories.StockRepository)
	loadedUser, err := stockRepo.GetByID(stockID)
	if err != nil {
		c.JSON(err.HttpStatus(), models.Envelope{Data: "", Message: err.Error()})
		return
	}
	if err := stockRepo.DeleteByObject(loadedUser); err != nil {
		c.JSON(err.HttpStatus(), models.Envelope{Data: "", Message: err.Error()})
		return
	}
	c.JSON(204, nil)
}
func stockGet(c *gin.Context) {
	stockID := c.Param("id")
	stockRepo := c.MustGet("stockRepo").(*repositories.StockRepository)
	loadedUser, err := stockRepo.GetByID(stockID)
	if err != nil {
		c.JSON(err.HttpStatus(), models.Envelope{Data: "", Message: err.Error()})
		return
	}
	env := models.Envelope{Data: loadedUser, Message: ""}
	c.JSON(200, env)
}
