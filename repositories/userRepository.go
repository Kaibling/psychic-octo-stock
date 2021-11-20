package repositories

import (
	"github.com/Kaibling/psychic-octo-stock/lib/apierrors"
	"github.com/Kaibling/psychic-octo-stock/lib/database"
	"github.com/Kaibling/psychic-octo-stock/models"
	"github.com/lucsky/cuid"
)

type UserRepository struct {
	db database.DBConnector
}

func NewUserRepository(dbConn database.DBConnector) *UserRepository {
	return &UserRepository{db: dbConn}
}

func (s *UserRepository) Add(user *models.User) apierrors.ApiError {
	user.ID = cuid.New()
	if err := s.db.Add(&user); err != nil {
		return err
	}
	return nil
}

func (s *UserRepository) GetByID(id string) (*models.User, apierrors.ApiError) {
	var object models.User

	if err := s.db.FindByID(&object, id, models.UserSelect); err != nil {
		return nil, err
	}
	return &object, nil

}
func (s *UserRepository) UpdateWithObject(user *models.User) apierrors.ApiError {
	if err := s.db.UpdateByObject(user); err != nil {
		return err
	}
	return nil
}

func (s *UserRepository) UpdateWithMap(data map[string]interface{}) apierrors.ApiError {
	if err := s.db.UpdateByMap(models.User{}, data); err != nil {
		return err
	}
	return nil
}

func (s *UserRepository) GetAll() ([]*models.User, apierrors.ApiError) {
	var userList []*models.User
	if err := s.db.GetAll(&userList, models.UserSelect); err != nil {
		return nil, err
	}
	return userList, nil
}

func (s *UserRepository) DeleteByID(data *models.User) apierrors.ApiError {
	if err := s.db.DeleteByID(data); err != nil {
		return err
	}
	return nil
}
