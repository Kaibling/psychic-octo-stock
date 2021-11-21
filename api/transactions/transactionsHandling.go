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

func ChangeStatus(transactionID string, status string) apierrors.ApiError {
	if !models.IsTransactionStatus(status) {
		return apierrors.NewClientError(errors.New("status '" + status + "' not valid"))
	}
	transactionRepo := repositories.TransactionRepo
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
			if !enoughStocks(loadedTransaction.SellerID, loadedTransaction.ID) {
				return apierrors.NewClientError(errors.New("seller does not have enough stocks"))
			}
			//set Status to ACTIVE
			updateData := map[string]interface{}{"ID": loadedTransaction.ID, "Status": "ACTIVE"}
			if err := transactionRepo.UpdateWithMap(updateData); err != nil {
				return apierrors.NewClientError(err)
			}
			log.Println(loadedTransaction.ID)
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
			updateData := map[string]interface{}{"ID": loadedTransaction.ID, "Status": "ACTIVE"}
			if err := transactionRepo.UpdateWithMap(updateData); err != nil {
				return apierrors.NewClientError(err)
			}
		}
	case "CLOSED":
		//cloesd, when status ACTIVE and :
		if loadedTransaction.Status != "ACTIVE" {
			return apierrors.NewClientError(errors.New("status '" + loadedTransaction.Status + "' wrong status for transition to CLOSED"))
		}

		if loadedTransaction.BuyerID == "" {
			return apierrors.NewClientError(errors.New("buyerID not set"))
		}
		if loadedTransaction.SellerID == "" {
			return apierrors.NewClientError(errors.New("sellerID not set"))
		}
		if !enoughStocks(loadedTransaction.SellerID, loadedTransaction.ID) {
			return apierrors.NewClientError(errors.New("seller does not have enough stocks"))
		}
		if !enoughFunds(loadedTransaction.BuyerID, loadedTransaction.ID) {
			return apierrors.NewClientError(errors.New("buyer has not enough funds"))
		}
		//execute transaction
		if err := executeTransaction(transactionID); err != nil {
			return apierrors.NewGeneralError(err)
		}

		//set Status to CLOSED
		updateData := map[string]interface{}{"ID": loadedTransaction.ID, "Status": "CLOSED"}
		if err := transactionRepo.UpdateWithMap(updateData); err != nil {
			return apierrors.NewClientError(err)
		}
	case "CANCELLED":
		//CANCELLED, when status ACTIVE
		if loadedTransaction.Status != "ACTIVE" {
			return apierrors.NewClientError(errors.New("status '" + loadedTransaction.Status + "' wrong status for transition to CANCELLED"))
		}
		updateData := map[string]interface{}{"ID": loadedTransaction.ID, "Status": "CANCELLED"}
		if err := transactionRepo.UpdateWithMap(updateData); err != nil {
			return apierrors.NewClientError(err)
		}

	case "PENDING":
		//CANCELLED, when status ACTIVE
		if loadedTransaction.Status != "ACTIVE" {
			return apierrors.NewClientError(errors.New("status '" + loadedTransaction.Status + "' wrong status for transition to PENDING"))
		}
		updateData := map[string]interface{}{"ID": loadedTransaction.ID, "Status": "PENDING"}
		if err := transactionRepo.UpdateWithMap(updateData); err != nil {
			return apierrors.NewClientError(err)
		}
	default:
		return apierrors.NewGeneralError(errors.New("status '" + status + "' not in switch statement"))
	}
	return nil
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
		log.Println("transaction cost not calculatable: " + err.Error())
		return false
	}
	return transactionCost <= userFunds
}

func enoughStocks(userID string, transactionID string) bool {
	transactionRepo := repositories.TransactionRepo
	transactionsData, err := transactionRepo.GetByID(transactionID)
	if err != nil {
		log.Println("stock fetch failed: " + err.Error())
		return false
	}
	stockRepo := repositories.StockRepo
	stockToUserData, err := stockRepo.GetStockPerUser(transactionsData.StockID,userID)
	if err != nil {
		log.Println("user has not stocks: " + err.Error())
		return false
	}
	return transactionsData.Quantity <= stockToUserData.Quantity
}

func executeTransaction(transactionID string) error {
	transactionRepo := repositories.TransactionRepo
	loadedTransaction, err := transactionRepo.GetByID(transactionID)
	if err != nil {
		return err
	}
	transactionCost, _ := transactionRepo.TransactionCostsbyID(transactionID)
	if err != nil {
		log.Println("funds fetch failed: " + err.Error())
	}
	userRepo := repositories.UserRepo
	seller, _ := userRepo.GetByID(loadedTransaction.SellerID)

	atomicExecutionArray := []interface{}{}
	updateSellerUserData := map[string]interface{}{"ID": loadedTransaction.SellerID, "Funds": seller.Funds + transactionCost}
	atomicExecutionArray = append(atomicExecutionArray, []interface{}{models.User{}, updateSellerUserData})

	stockRepo := repositories.StockRepo
	stockToUserDataSeller, err := stockRepo.GetStockPerUser(loadedTransaction.SellerID, loadedTransaction.StockID)
	if err != nil {
		log.Println("funds fetch failed: " + err.Error())
	}

	updateSellerStockData := map[string]interface{}{"StockID": loadedTransaction.StockID, "UserID": loadedTransaction.SellerID, "Quantity": stockToUserDataSeller.Quantity - loadedTransaction.Quantity}
	atomicExecutionArray = append(atomicExecutionArray, []interface{}{models.StockToUser{}, updateSellerStockData})

	updateBuyerUserData := map[string]interface{}{"ID": loadedTransaction.BuyerID, "Funds": seller.Funds - transactionCost}
	atomicExecutionArray = append(atomicExecutionArray, []interface{}{models.User{}, updateBuyerUserData})
	//update or insert
	stockToUserDataBuyer, err := stockRepo.GetStockPerUser(loadedTransaction.SellerID, loadedTransaction.StockID)
	if err != nil {
		log.Println("funds fetch failed: " + err.Error())
	}

	if stockToUserDataBuyer == nil {
		//Buyer has not Stocks. Insert instead of update
		insertBuyerStockData := map[string]interface{}{"StockID": loadedTransaction.StockID, "UserID": loadedTransaction.BuyerID, "Quantity": loadedTransaction.Quantity}
		stockRepo.AddStockToUser(loadedTransaction.StockID, loadedTransaction.BuyerID, 0)
		atomicExecutionArray = append(atomicExecutionArray, []interface{}{models.StockToUser{}, insertBuyerStockData})
	} else {
		updateBuyerStockData := map[string]interface{}{"StockID": loadedTransaction.StockID, "UserID": loadedTransaction.BuyerID, "Quantity": stockToUserDataBuyer.Quantity + loadedTransaction.Quantity}
		atomicExecutionArray = append(atomicExecutionArray, []interface{}{models.StockToUser{}, updateBuyerStockData})
	}

	if err := transactionRepo.ExecuteTransaction(atomicExecutionArray); err != nil {
		return apierrors.NewGeneralError(err)
	}
	//ExecuteTransaction(userData1 , userData2 map[string]interface{})
	//loadedTransaction
	return nil
}
