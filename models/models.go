package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       string  `gorm:"primaryKey;autoIncrement:false;not null"`
	Username string  `gorm:"unique;not null" json:"username"`
	Password string  `gorm:"not null" json:"password"`
	Email    string  `gorm:"unique" json:"email"`
	Address  string  `json:"address"`
	Funds    float64 `gorm:"default:0" json:"funds"`
}

var UserSelect = []string{"ID", " Username", " Email", " Address", "Funds"}

type Stock struct {
	gorm.Model
	ID       string `gorm:"primaryKey;autoIncrement:false;not null"`
	Name     string `gorm:"not null;unique" json:"name"`
	Quantity int    `gorm:"not null" json:"quantity"`
}

var StockSelect = []string{"ID", " Name", "Quantity"}

type StockToUser struct {
	gorm.Model
	ID       string `gorm:"primaryKey;autoIncrement:false;not null"`
	StockID  string `gorm:"foreignkey:StockID" json:"stock_id"`
	UserID   string `gorm:"foreignkey:UserID" json:"user_id"`
	Quantity int    `gorm:"not null" json:"quantity"`
}

type Transaction struct {
	gorm.Model
	ID       string  `gorm:"primaryKey;autoIncrement:false;not null"`
	SellerID string  `gorm:"foreignkey:userID" json:"seller_id"`
	BuyerID  string  `gorm:"foreignkey:userID" json:"buyer_id"`
	StockID  string  `gorm:"foreignkey:stockID;not null" json:"stock_id"`
	Quantity int     `gorm:"not null" json:"quantity"`
	Price    float64 `gorm:"not null" json:"price"`
	Type     string  `gorm:"not null" json:"type"`
	Status   string  `gorm:"not null;default:PENDING" json:"status"`
	Comment  string  `gorm:"default:initiated" json:"comment"`
}

var transactionTypes = []string{"SELL", "BUY"}
var transactionStatus = []string{"PENDING", "ACTIVE", "CLOSED", "CANCELLED"}
var TransactionSelect = []string{"ID", "seller_id", "buyer_id", "stock_id", " Quantity", " Type", "Price", "Status"}

func IsTransactionsType(data string) bool {
	return contains(transactionTypes, data)
}

func IsTransactionStatus(data string) bool {
	return contains(transactionStatus, data)

}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

type Token struct {
	gorm.Model
	ID      string `gorm:"primaryKey;autoIncrement:false;not null"`
	Active  bool   `gorm:"active"`
	Comment string `gorm:"comment"`
	UserID  string `gorm:"foreignkey:userID;not null" json:"user_id"`
}

var TokenSelect = []string{"ID", "user_id", "comment", "active"}
