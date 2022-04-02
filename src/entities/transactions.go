package entities

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	WalletID uint            `json:"transaction_id"`
	Amount   decimal.Decimal `json:"amount"`
}
