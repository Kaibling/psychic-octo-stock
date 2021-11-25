package transactions

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Kaibling/psychic-octo-stock/lib/apierrors"
	"github.com/Kaibling/psychic-octo-stock/lib/transmission"
	"github.com/Kaibling/psychic-octo-stock/lib/utility"
	"github.com/Kaibling/psychic-octo-stock/models"
	"github.com/Kaibling/psychic-octo-stock/repositories"
	"github.com/go-chi/chi"
)

func transactionPost(w http.ResponseWriter, r *http.Request) {
	response := transmission.GetOrCreateResponse(w, r)
	var newTransaction models.Transaction
	erra := json.NewDecoder(r.Body).Decode(&newTransaction)
	if erra != nil {
		response.Send("", "post data not parsable", http.StatusUnprocessableEntity)
		return

	}
	if !models.IsTransactionsType(newTransaction.Type) {
		err := apierrors.NewClientError(errors.New("transactiontype is neither 'BUY' or 'SELL'"))
		response.Send("", err.Error(), err.HttpStatus())
		return
	}
	//todo check if seller has enoigh stock
	transactionRepo := utility.GetContext("transactionRepo", r).(*repositories.TransactionRepository)
	if err := transactionRepo.Add(&newTransaction); err != nil {
		response.Send("", err.Error(), err.HttpStatus())
		return
	}
	//todo proper return schema
	response.Send(newTransaction, "", http.StatusCreated)
}
func transactionsGet(w http.ResponseWriter, r *http.Request) {
	response := transmission.GetOrCreateResponse(w, r)
	transactionRepo := utility.GetContext("transactionRepo", r).(*repositories.TransactionRepository)
	transactionList, err := transactionRepo.GetAll()
	if err != nil {
		response.Send("", err.Error(), err.HttpStatus())
		return
	}

	response.Send(transactionList, "", http.StatusOK)
}

func transactionDelete(w http.ResponseWriter, r *http.Request) {
	response := transmission.GetOrCreateResponse(w, r)
	transactionID := chi.URLParam(r, "id")
	transactionRepo := utility.GetContext("transactionRepo", r).(*repositories.TransactionRepository)
	loadedUser, err := transactionRepo.GetByID(transactionID)
	if err != nil {
		response.Send("", err.Error(), err.HttpStatus())
		return
	}
	if err := transactionRepo.DeleteByObject(loadedUser); err != nil {
		response.Send("", err.Error(), err.HttpStatus())
		return
	}
	response.Send("", "", http.StatusNoContent)
}

func transactionGet(w http.ResponseWriter, r *http.Request) {
	response := transmission.GetOrCreateResponse(w, r)
	transactionID := chi.URLParam(r, "id")
	transactionRepo := utility.GetContext("transactionRepo", r).(*repositories.TransactionRepository)
	loadedUser, err := transactionRepo.GetByID(transactionID)
	if err != nil {
		response.Send("", err.Error(), err.HttpStatus())
		return
	}
	response.Send(loadedUser, "", http.StatusOK)
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
	stockToUserData, err := stockRepo.GetStockPerUser(transactionsData.StockID, userID)
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
		return err
	}
	userRepo := repositories.UserRepo
	seller, _ := userRepo.GetByID(loadedTransaction.SellerID)

	atomicExecutionArray := []interface{}{}
	updateSellerUserData := map[string]interface{}{"funds": seller.Funds + transactionCost}
	var updateSellerUserQuery interface{} = "id = ?"
	atomicExecutionArray = append(atomicExecutionArray, []interface{}{models.User{}, updateSellerUserData, updateSellerUserQuery, []interface{}{loadedTransaction.SellerID}})
	stockRepo := repositories.StockRepo
	stockToUserDataSeller, err := stockRepo.GetStockPerUser(loadedTransaction.StockID, loadedTransaction.SellerID)
	if err != nil {
		log.Println("Seller has not enough Stocks: " + err.Error())
		return err
	}

	updateSellerStockData := map[string]interface{}{"Quantity": stockToUserDataSeller.Quantity - loadedTransaction.Quantity}
	var updateSellerStockQuery interface{} = "stock_id = ? and user_id = ?"
	atomicExecutionArray = append(atomicExecutionArray, []interface{}{models.StockToUser{}, updateSellerStockData, updateSellerStockQuery, []interface{}{loadedTransaction.StockID, loadedTransaction.SellerID}})

	updateBuyerUserData := map[string]interface{}{"funds": seller.Funds - transactionCost}
	var updateBuyerUserQuery interface{} = "id = ?"
	atomicExecutionArray = append(atomicExecutionArray, []interface{}{models.User{}, updateBuyerUserData, updateBuyerUserQuery, []interface{}{loadedTransaction.BuyerID}})
	//update or insert
	stockToUserDataBuyer, err := stockRepo.GetStockPerUser(loadedTransaction.StockID, loadedTransaction.BuyerID)
	if err != nil {
		log.Println("Stocks record for Buyer: " + err.Error())
		//Buyer has not Stocks. Insert instead of update
		err := stockRepo.AddStockToUser(loadedTransaction.StockID, loadedTransaction.BuyerID, 0)
		if err != nil {
			log.Println("creation of stockToUser record failed:  " + err.Error())
			return err
		}

		insertBuyerStockData := map[string]interface{}{"Quantity": loadedTransaction.Quantity}
		var insertBuyerStockQuery interface{} = "stock_id = ? and user_id = ?"
		atomicExecutionArray = append(atomicExecutionArray, []interface{}{models.StockToUser{}, insertBuyerStockData, insertBuyerStockQuery, []interface{}{loadedTransaction.StockID, loadedTransaction.BuyerID}})
	} else {
		updateBuyerStockData := map[string]interface{}{"Quantity": stockToUserDataBuyer.Quantity + loadedTransaction.Quantity}
		var updateBuyerStockQuery interface{} = "stock_id = ? and user_id = ?"
		atomicExecutionArray = append(atomicExecutionArray, []interface{}{models.StockToUser{}, updateBuyerStockData, updateBuyerStockQuery, []interface{}{loadedTransaction.StockID, loadedTransaction.BuyerID}})
	}

	if err := transactionRepo.ExecuteTransaction(atomicExecutionArray); err != nil {
		return err
	}
	//ExecuteTransaction(userData1 , userData2 map[string]interface{})
	//loadedTransaction
	return nil
}
