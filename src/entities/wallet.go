package entities

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Wallet struct {
	gorm.Model
	Balance decimal.Decimal `json:"balance"`
	Active  bool            `json:"active"`
}
