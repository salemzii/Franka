package entities

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Wallet struct {
	gorm.Model
	UserID       uint            `json:"user_id"`
	Balance      decimal.Decimal `json:"balance"`
	Active       bool            `json:"active" gorm:"default:false"`
	Transactions []Transaction
}
