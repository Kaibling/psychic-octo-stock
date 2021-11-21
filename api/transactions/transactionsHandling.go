package transactions

import (
	"errors"

	"github.com/Kaibling/psychic-octo-stock/lib/apierrors"
	"github.com/Kaibling/psychic-octo-stock/models"
	"github.com/Kaibling/psychic-octo-stock/repositories"
	"github.com/gin-gonic/gin"
)

func transactionPost(c *gin.Context) {
	var newTransaction models.Transaction
	c.BindJSON(&newTransaction)
	if !models.IsTransactionsType(newTransaction.Type) {
		err := apierrors.NewClientError(errors.New("transactiontype is neither 'BUY' or 'SELL'"))
		c.JSON(err.HttpStatus(), models.Envelope{Data: "", Message: err.Error()})
		return
	}
	//todo check if seller has enoigh stock
	transactionRepo := c.MustGet("transactionRepo").(*repositories.TransactionRepository)
	if err := transactionRepo.Add(&newTransaction); err != nil {
		c.JSON(err.HttpStatus(), models.Envelope{Data: "", Message: err.Error()})
		return
	}
	//todo proper return schema
	env := models.Envelope{Data: newTransaction, Message: ""}
	c.JSON(201, env)
}
func transactionsGet(c *gin.Context) {
	transactionRepo := c.MustGet("transactionRepo").(*repositories.TransactionRepository)
	transactionList, err := transactionRepo.GetAll()
	if err != nil {
		c.JSON(err.HttpStatus(), models.Envelope{Data: "", Message: err.Error()})
		return
	}
	env := models.Envelope{Data: transactionList, Message: ""}
	c.JSON(200, env)
}

func transactionDelete(c *gin.Context) {
	transactionID := c.Param("id")
	transactionRepo := c.MustGet("transactionRepo").(*repositories.TransactionRepository)
	loadedUser, err := transactionRepo.GetByID(transactionID)
	if err != nil {
		c.JSON(err.HttpStatus(), models.Envelope{Data: "", Message: err.Error()})
		return
	}
	if err := transactionRepo.DeleteByID(loadedUser); err != nil {
		c.JSON(err.HttpStatus(), models.Envelope{Data: "", Message: err.Error()})
		return
	}
	c.JSON(204, nil)
}
func transactionGet(c *gin.Context) {
	transactionID := c.Param("id")
	transactionRepo := c.MustGet("transactionRepo").(*repositories.TransactionRepository)
	loadedUser, err := transactionRepo.GetByID(transactionID)
	if err != nil {
		c.JSON(err.HttpStatus(), models.Envelope{Data: "", Message: err.Error()})
		return
	}
	env := models.Envelope{Data: loadedUser, Message: ""}
	c.JSON(200, env)
}
