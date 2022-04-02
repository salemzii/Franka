package entities

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	Wallet   Wallet
	WalletId uint
	Amount   decimal.Decimal
}
