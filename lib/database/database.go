package database

import (
	"encoding/json"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DBConnector interface {
	Connect()
	Migrate(interface{})
	Add(interface{})
	FindByID(object interface{}, id string) interface{}
	UpdateByObject(data interface{})
}

type GormConnector struct {
	url       string
	connector *gorm.DB
}

func NewDatabaseConnector(url string) *GormConnector {
	return &GormConnector{url: url}
}

func (s *GormConnector) Connect() {
	db, err := gorm.Open(sqlite.Open(s.url), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	s.connector = db
}

func (s *GormConnector) Migrate(object interface{}) {
	s.connector.AutoMigrate(&object)
}

func (s *GormConnector) Add(object interface{}) {
	b, _ := json.MarshalIndent(object, "", "  ")
	fmt.Println(string(b))
	s.connector.Create(object)
}
func (s *GormConnector) FindByID(object interface{}, id string) interface{} {
	s.connector.First(&object, "id = ?", id)
	if object == nil {
		fmt.Println("nothing found")
	}
	return object
}
func (s *GormConnector) UpdateByObject(data interface{}) {
	s.connector.Save(data)
}
