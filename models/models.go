package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       string `gorm:"primaryKey"`
	Username string
	Password string
	Email    string
	Address  string
}

type Stock struct {
	gorm.Model
	ID   string `gorm:"primaryKey"`
	Name string
}
