package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       string `gorm:"primaryKey;autoIncrement:false"`
	Username string `gorm:"unique" json:"username"`
	Password string `json:"password"`
	Email    string `gorm:"unique" json:"email"`
	Address  string `json:"address"`
}

var UserSelect = []string{"ID", " Username", " Email", " Address"}

type Stock struct {
	gorm.Model
	ID   string `gorm:"primaryKey;autoIncrement:false"`
	Name string `gorm:"unique"`
}

type Envelope struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}
