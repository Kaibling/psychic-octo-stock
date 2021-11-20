package database

import (
	"errors"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DBConnector interface {
	Connect() error
	Migrate(interface{}) error
	Add(interface{}) error
	FindByID(object interface{}, id string) error
	UpdateByObject(data interface{}) error
}

type GormConnector struct {
	url       string
	connector *gorm.DB
}

func NewDatabaseConnector(url string) *GormConnector {
	return &GormConnector{url: url}
}

func (s *GormConnector) Connect() error {
	db, err := gorm.Open(sqlite.Open(s.url), &gorm.Config{})
	if err != nil {
		return err
		//panic("failed to connect database")
	}
	s.connector = db
	return nil
}

func (s *GormConnector) Migrate(object interface{}) error {
	s.connector.AutoMigrate(&object)
	return nil
}

func (s *GormConnector) Add(object interface{}) error {
	if dbc := s.connector.Create(object); dbc.Error != nil {
		return dbc.Error
	}
	return nil
}
func (s *GormConnector) FindByID(object interface{}, id string) error {
	s.connector.First(&object, "id = ?", id)
	if object == nil {
		return errors.New("nothing found")
	}
	return nil
}
func (s *GormConnector) UpdateByObject(data interface{}) error {
	s.connector.Save(data)
	return nil
}
