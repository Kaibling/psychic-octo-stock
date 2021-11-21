package repositories

import (
	"github.com/Kaibling/psychic-octo-stock/lib/apierrors"
	"github.com/Kaibling/psychic-octo-stock/lib/database"
	"github.com/Kaibling/psychic-octo-stock/models"
	"github.com/lucsky/cuid"
)

var UserRepo *UserRepository

type UserRepository struct {
	db database.DBConnector
}

func SetUserRepo(repo *UserRepository) {
	UserRepo = repo
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

func (s *UserRepository) DeleteByObject(data *models.User) apierrors.ApiError {
	if err := s.db.DeleteByID(data); err != nil {
		return err
	}
	return nil
}

func (s *UserRepository) FundsByID(id string) (float64, apierrors.ApiError) {
	var user *models.User
	selectString := []string{"funds"}
	if err := s.db.GetData(&user, selectString, id); err != nil {
		return 0, err
	}

	return user.Funds, nil
}
