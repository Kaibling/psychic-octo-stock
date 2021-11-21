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

var UserSelect = []string{"ID", " Username", " Email", " Address"}

type Stock struct {
	gorm.Model
	ID       string `gorm:"primaryKey;autoIncrement:false;not null"`
	Name     string `gorm:"not null;unique" json:"name"`
	Quantity int    `gorm:"not null" json:"quantity"`
}

var StockSelect = []string{"ID", " Name", "Quantity"}

type StockToUser struct {
	gorm.Model
	//ID       string `gorm:"primaryKey;autoIncrement:false;not null"`
	StockID  string `gorm:"foreignkey:StockID;primaryKey"`
	UserID   string `gorm:"foreignkey:UserID;primaryKey"`
	Quantity int    `gorm:"not null" json:"quantity"`
}

type Transaction struct {
	gorm.Model
	ID       string  `gorm:"primaryKey;autoIncrement:false;not null"`
	SellerID string  `gorm:"foreignkey:userID" json:"sellerID"`
	BuyerID  string  `gorm:"foreignkey:userID" json:"buyerID"`
	StockID  string  `gorm:"foreignkey:stockID;not null" json:"stockID"`
	Quantity int     `gorm:"not null" json:"quantity"`
	Price    float64 `gorm:"not null" json:"price"`
	Type     string  `gorm:"not null" json:"type"`
	Status   string  `gorm:"not null;default:PENDING" json:"status"`
	Comment  string  `gorm:"default:initiated" json:"comment"`
}

var transactionTypes = []string{"SELL", "BUY"}
var transactionStatus = []string{"PENDING", "ACTIVE", "CLOSED", "CANCELLED"}
var TransactionSelect = []string{"ID", " seller_id", "buyer_id", "stock_id", " Quantity", " Type"}

type Envelope struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

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
