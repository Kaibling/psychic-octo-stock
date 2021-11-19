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

func (s *UserRepository) AddUser(user *models.User) {
	user.ID = cuid.New()
	s.db.Add(&user)
}

func (s *UserRepository) GetUserByID(id string) *models.User {
	var object models.User
	s.db.FindByID(&object, id)
	return &object

}
func (s *UserRepository) UpdateObject(user *models.User) {
	s.db.UpdateByObject(user)
}
