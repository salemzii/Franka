package entities

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	WalletID uint            `json:"transaction_id"`
	Type     string          `json:"type"`
	Amount   decimal.Decimal `json:"amount"`
}
