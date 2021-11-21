package repositories

import (
	"github.com/Kaibling/psychic-octo-stock/lib/apierrors"
	"github.com/Kaibling/psychic-octo-stock/lib/database"
	"github.com/Kaibling/psychic-octo-stock/models"
	"github.com/lucsky/cuid"
)

type TransactionRepository struct {
	db database.DBConnector
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

func (s *TransactionRepository) DeleteByID(data *models.Transaction) apierrors.ApiError {
	if err := s.db.DeleteByID(data); err != nil {
		return err
	}
	return nil
}