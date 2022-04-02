package entities

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Wallet struct {
	gorm.Model
	Balance      decimal.Decimal `json:"balance"`
	Active       bool            `json:"active" gorm:"default:false"`
	WalletId     string          `json:"wallet_id"`
	Transactions []Transaction   `json:"transactions" gorm:"foreignKey:UserRefer"`
}
