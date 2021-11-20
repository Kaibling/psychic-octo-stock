package repositories

import (
	"log"

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

func (s *UserRepository) AddUser(user *models.User) error {
	user.ID = cuid.New()
	if err := s.db.Add(&user); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *UserRepository) GetUserByID(id string) *models.User {
	var object models.User
	s.db.FindByID(&object, id, models.UserSelect)
	return &object

}
func (s *UserRepository) UpdateWithObject(user *models.User) error {
	s.db.UpdateByObject(user)
	return nil
}

func (s *UserRepository) UpdateWithMap(data map[string]interface{}) error {
	s.db.UpdateByMap(models.User{}, data)
	return nil
}

func (s *UserRepository) GetAllUser() []*models.User {
	var userList []*models.User
	s.db.GetAll(&userList, models.UserSelect)
	return userList
}
