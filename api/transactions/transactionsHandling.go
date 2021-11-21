package transactions

import (
	"errors"
	"log"

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
	if err := transactionRepo.DeleteByObject(loadedUser); err != nil {
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

func changeStatus(transactionID string, status string, transactionRepo *repositories.TransactionRepository) apierrors.ApiError {
	if !models.IsTransactionStatus(status) {
		return apierrors.NewClientError(errors.New("status '" + status + "' not valid"))
	}
	loadedTransaction, _ := transactionRepo.GetByID(transactionID)
	switch supportedStatus := status; supportedStatus {
	case "ACTIVE":
		//active, when status PENDING and :
		if loadedTransaction.Status != "PENDING" {
			return apierrors.NewClientError(errors.New("status '" + loadedTransaction.Status + "' wrong status for transition to ACTIVE"))
		}
		if loadedTransaction.Type == "SELL" {
			//type SELL and sellerID not "" and buyerid "" and seller has enough stock
			if loadedTransaction.SellerID == "" {
				return apierrors.NewClientError(errors.New("sellerID not set"))
			}
			if loadedTransaction.BuyerID != "" {
				return apierrors.NewClientError(errors.New("buyerID already set"))
			}
			// if enoughFunds(userID string, transactionID string) {

			// }
		}
		if loadedTransaction.Type == "BUY" {
			//type BUY and buyerID not "" and sellerid "" and buyer has enough funds
			if loadedTransaction.SellerID != "" {
				return apierrors.NewClientError(errors.New("sellerID already set"))
			}
			if loadedTransaction.BuyerID == "" {
				return apierrors.NewClientError(errors.New("buyerID not set"))
			}
			if !enoughFunds(loadedTransaction.BuyerID, loadedTransaction.ID) {
				return apierrors.NewClientError(errors.New("buyer has not enough funds"))
			}
		}
	case "CLOSED":
		//cloesd, when status ACTIVE and :
		//type SELL and buyerID not "" and sellerid "" and buyer has enough funds
		//type BUY  and sellerID not "" and buyerid "" and seller has enough stock
	case "CANCELLED":
		//CANCELLED, when status ACTIVE

	case "PENDING":
		//CANCELLED, when status ACTIVE
		//add comment
	default:
		return  apierrors.NewGeneralError(errors.New("status '"+ status +  "' not in switch statement"))
	}
	return  nil
}

func enoughFunds(userID string, transactionID string) bool {
	// get available funds of user
	userRepo := repositories.UserRepo
	userFunds, err := userRepo.FundsByID(userID)
	if err != nil {
		log.Println("funds fetch failed: " + err.Error())
		return false
	}
	transactionRepo := repositories.TransactionRepo
	transactionCost, _ := transactionRepo.TransactionCostsbyID(transactionID)
	if err != nil {
		log.Println("funds fetch failed: " + err.Error())
		return false
	}
	return transactionCost <= userFunds
}


func enoughStocks(userID string, transactionID string) bool {
	// get available funds of user
	userRepo := repositories.UserRepo
	userFunds, err := userRepo.FundsByID(userID)
	if err != nil {
		log.Println("funds fetch failed: " + err.Error())
		return false
	}
	transactionRepo := repositories.TransactionRepo
	transactionCost, _ := transactionRepo.TransactionCostsbyID(transactionID)
	if err != nil {
		log.Println("funds fetch failed: " + err.Error())
		return false
	}
	return transactionCost <= userFunds
}