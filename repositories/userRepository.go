package repositories

import (
	"github.com/Kaibling/psychic-octo-stock/lib/apierrors"
	"github.com/Kaibling/psychic-octo-stock/lib/database"
	"github.com/Kaibling/psychic-octo-stock/lib/utility"
	"github.com/Kaibling/psychic-octo-stock/models"
	"github.com/Kaibling/psychic-octo-stock/modules"
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
	user.Password = utility.HashPassword(user.Password)
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
	if err := s.db.DeleteByObject(data); err != nil {
		return err
	}
	return nil
}

func (s *UserRepository) FundsByID(id string) (*models.MonetaryUnit, apierrors.ApiError) {
	var user *models.User
	selectString := []string{"funds", "currency"}
	if err := s.db.GetData(&user, selectString, id); err != nil {
		return nil, err
	}

	return &models.MonetaryUnit{Amount: user.Funds, Currency: user.Currency}, nil
}

func (s *UserRepository) GetPWByName(userName string) (string, apierrors.ApiError) {
	var object models.User

	query := "username = ?"
	querData := []interface{}{userName}
	if err := s.db.FindByWhere(&object, query, querData); err != nil {
		return "", err
	}
	return object.Password, nil
}
func (s *UserRepository) GetyName(userName string) (*models.User, apierrors.ApiError) {
	var object models.User

	query := "username = ?"
	querData := []interface{}{userName}
	if err := s.db.FindByWhere(&object, query, querData); err != nil {
		return nil, err
	}
	return &object, nil
}

func (s *UserRepository) AddFunds(userID string, mu models.MonetaryUnit) apierrors.ApiError {
	currentFunds, err := s.FundsByID(userID)
	if err != nil {
		return err
	}
	modules.CCM.AddAndConvertFunds(currentFunds, mu)
	updateUser := map[string]interface{}{"ID": userID, "Funds": currentFunds.Amount}

	if err := s.db.UpdateByMap(models.User{}, updateUser); err != nil {
		return err
	}
	return nil
}
