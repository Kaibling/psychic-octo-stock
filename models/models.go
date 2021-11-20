package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       string `gorm:"primaryKey"`
	Username string `json:"username"` //`gorm:"unique" json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"` //`gorm:"unique" json:"email"`
	Address  string `json:"address"`
}

type Stock struct {
	gorm.Model
	ID   string `gorm:"primaryKey"`
	Name string `gorm:"unique"`
}

type Envelope struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}
