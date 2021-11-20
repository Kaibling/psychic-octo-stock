package repositories

import (
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
		//log.Println(user.ID)
		return err
	}
	return nil
}

func (s *UserRepository) GetUserByID(id string) *models.User {
	var object models.User
	s.db.FindByID(&object, id)
	return &object

}
func (s *UserRepository) UpdateObject(user *models.User) error {
	s.db.UpdateByObject(user)
	return nil
}
