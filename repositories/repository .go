package repositories

import (
	"github.com/Kaibling/psychic-octo-stock/lib/apierrors"
	"github.com/Kaibling/psychic-octo-stock/lib/database"
)

type Repository struct {
	db           database.DBConnector
	dbModel      interface{}
	selectString []string
}

func NewRepository(dbConn database.DBConnector, dbModel interface{}, selectString []string) *Repository {
	return &Repository{db: dbConn, dbModel: dbModel, selectString: selectString}
}

func (s *Repository) Add(user *interface{}) apierrors.ApiError {
	if err := s.db.Add(&user); err != nil {
		return err
	}
	return nil
}

func (s *Repository) GetByID(object interface{}, id string) apierrors.ApiError {

	if err := s.db.FindByID(&object, id, s.selectString); err != nil {
		return err
	}
	return nil

}
func (s *Repository) UpdateWithObject(object interface{}) apierrors.ApiError {
	if err := s.db.UpdateByObject(object); err != nil {
		return err
	}
	return nil
}

func (s *Repository) UpdateWithMap(data map[string]interface{}) apierrors.ApiError {
	if err := s.db.UpdateByMap(s.dbModel, data); err != nil {
		return err
	}
	return nil
}

func (s *Repository) GetAll(list []*interface{}) apierrors.ApiError {
	if err := s.db.GetAll(&list, s.selectString); err != nil {
		return err
	}
	return nil
}

func (s *Repository) DeleteByObject(data interface{}) apierrors.ApiError {
	if err := s.db.DeleteByObject(data); err != nil {
		return err
	}
	return nil
}

func (s *Repository) FundsByID(data interface{}) apierrors.ApiError {
	if err := s.db.DeleteByObject(data); err != nil {
		return err
	}
	return nil
}
