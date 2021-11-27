package repositories

import (
	"github.com/Kaibling/psychic-octo-stock/lib/apierrors"
	"github.com/Kaibling/psychic-octo-stock/lib/database"
	"github.com/Kaibling/psychic-octo-stock/models"
	"github.com/lucsky/cuid"
)

var StockRepo *StockRepository

type StockRepository struct {
	db database.DBConnector
}

func SetStockRepo(repo *StockRepository) {
	StockRepo = repo
}

func NewStockRepository(dbConn database.DBConnector) *StockRepository {
	return &StockRepository{db: dbConn}
}

func (s *StockRepository) Add(stock *models.Stock) apierrors.ApiError {
	stock.ID = cuid.New()
	if err := s.db.Add(&stock); err != nil {
		return err
	}
	return nil
}

func (s *StockRepository) GetByID(id string) (*models.Stock, apierrors.ApiError) {
	var object models.Stock

	if err := s.db.FindByID(&object, id, models.StockSelect); err != nil {
		return nil, err
	}
	return &object, nil

}
func (s *StockRepository) UpdateWithObject(stock *models.Stock) apierrors.ApiError {
	if err := s.db.UpdateByObject(stock); err != nil {
		return err
	}
	return nil
}

func (s *StockRepository) UpdateWithMap(data map[string]interface{}) apierrors.ApiError {
	if err := s.db.UpdateByMap(models.Stock{}, data); err != nil {
		return err
	}
	return nil
}

func (s *StockRepository) GetAll() ([]*models.Stock, apierrors.ApiError) {
	var stockList []*models.Stock
	if err := s.db.GetAll(&stockList, models.StockSelect); err != nil {
		return nil, err
	}
	return stockList, nil
}

func (s *StockRepository) DeleteByObject(data *models.Stock) apierrors.ApiError {
	if err := s.db.DeleteByObject(data); err != nil {
		return err
	}
	return nil
}

func (s *StockRepository) AddStockToUser(stockID string, userID string, quantity int) apierrors.ApiError {
	model := &models.StockToUser{ID: cuid.New(), UserID: userID, StockID: stockID, Quantity: quantity}
	if err := s.db.Add(&model); err != nil {
		return err
	}
	return nil
}

func (s *StockRepository) GetStockPerUser(stockID string, userID string) (*models.StockToUser, apierrors.ApiError) {
	var object models.StockToUser

	query := "user_id = ? AND stock_id = ?"
	querData := []interface{}{userID, stockID}
	if err := s.db.FindByWhere(&object, query, querData); err != nil {
		return nil, err
	}
	return &object, nil
}
