package repositories

import (
	"github.com/Kaibling/psychic-octo-stock/lib/apierrors"
	"github.com/Kaibling/psychic-octo-stock/lib/database"
	"github.com/Kaibling/psychic-octo-stock/models"
	"github.com/lucsky/cuid"
)

var TransactionRepo *TransactionRepository

type TransactionRepository struct {
	db database.DBConnector
}

func SetTransactionRepo(repo *TransactionRepository) {
	TransactionRepo = repo
}

func NewTransactionRepository(dbConn database.DBConnector) *TransactionRepository {
	return &TransactionRepository{db: dbConn}
}

func (s *TransactionRepository) Add(Transaction *models.Transaction) apierrors.ApiError {
	Transaction.ID = cuid.New()
	if err := s.db.Add(&Transaction); err != nil {
		return err
	}
	return nil
}

func (s *TransactionRepository) GetByID(id string) (*models.Transaction, apierrors.ApiError) {
	var object models.Transaction

	if err := s.db.FindByID(&object, id, models.TransactionSelect); err != nil {
		return nil, err
	}
	return &object, nil

}
func (s *TransactionRepository) UpdateWithObject(Transaction *models.Transaction) apierrors.ApiError {
	if err := s.db.UpdateByObject(Transaction); err != nil {
		return err
	}
	return nil
}

func (s *TransactionRepository) UpdateWithMap(data map[string]interface{}) apierrors.ApiError {
	if err := s.db.UpdateByMap(models.Transaction{}, data); err != nil {
		return err
	}
	return nil
}

func (s *TransactionRepository) GetAll() ([]*models.Transaction, apierrors.ApiError) {
	var TransactionList []*models.Transaction
	if err := s.db.GetAll(&TransactionList, models.TransactionSelect); err != nil {
		return nil, err
	}
	return TransactionList, nil
}

func (s *TransactionRepository) DeleteByObject(data *models.Transaction) apierrors.ApiError {
	if err := s.db.DeleteByObject(data); err != nil {
		return err
	}
	return nil
}

func (s *TransactionRepository) TransactionCostsbyID(transactionID string) (*models.MonetaryUnit, apierrors.ApiError) {
	var transaction *models.Transaction
	selectString := []string{"price", "quantity", "currency"}
	if err := s.db.GetData(&transaction, selectString, transactionID); err != nil {
		return nil, err
	}
	combinedPrice := float64(transaction.Quantity) * transaction.Price
	return &models.MonetaryUnit{Amount: combinedPrice, Currency: transaction.Currency}, nil

}

func (s *TransactionRepository) ExecuteTransaction(Data []interface{}) apierrors.ApiError {
	if err := s.db.ExecuteTransaction(Data); err != nil {
		return err
	}
	return nil
}
