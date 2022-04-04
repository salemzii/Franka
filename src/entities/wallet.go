package entities

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
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

func GetWallet(c *gin.Context) {
	wallet_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
	}
	var wallet Wallet
	db.Transaction(func(tx *gorm.DB) error {
		tx.First(&wallet, wallet_id)
		return nil
	})

	c.JSON(200, gin.H{
		"id":           wallet.ID,
		"balance":      wallet.Balance,
		"user_id":      wallet.UserID,
		"transactions": wallet.Transactions,
	})

}

type TransactionStruct struct {
	Amount decimal.Decimal `json:"amount"`
}

func CreditWallet(c *gin.Context) {
	wallet_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
	}
	var wallet Wallet
	var creditAmount TransactionStruct
	//var transaction Transaction
	c.BindJSON(&creditAmount)

	db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", wallet_id).First(&wallet).Error; err != nil {
			c.JSON(400, gin.H{
				"error": "wallet with id " + strconv.Itoa(wallet_id) + " not found!",
			})
		} else {
			if creditAmount.Amount.IsPositive() {
				if err := tx.Where("id = ?", wallet_id).First(&wallet).Update("balance", wallet.Balance.Add(creditAmount.Amount)).Error; err != nil {
					c.JSON(200, gin.H{
						"error": "Error while processing transaction",
					})
				}
				if err := tx.Where("id = ?", wallet_id).First(&wallet).Update("transaction", Transaction{WalletID: wallet.ID, Type: "debit", Amount: creditAmount.Amount}).Error; err != nil {
					c.JSON(200, gin.H{
						"error": "Error while updating transaction",
					})
				}
				c.JSON(200, gin.H{
					"wallet_balance": wallet.Balance,
				})
			} else {
				//return status 400 if amount to be credited is negative.
				c.JSON(400, gin.H{
					"error": "Cannot use negative value " + creditAmount.Amount.String() + " for operation",
				})
			}
		}
		return nil
	})
}

func DebitWallet(c *gin.Context) {
	wallet_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
	}
	var wallet Wallet
	var debitAmount TransactionStruct
	//var transaction Transaction
	c.BindJSON(&debitAmount)

	db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", wallet_id).First(&wallet).Error; err != nil {
			c.JSON(400, gin.H{
				"error": "wallet with id " + strconv.Itoa(wallet_id) + " not found!",
			})
		} else {
			if debitAmount.Amount.GreaterThan(wallet.Balance) {
				c.JSON(400, gin.H{
					"error": "operation not allowed! cannot debit value greater than wallet balance",
				})
			} else {
				if debitAmount.Amount.IsPositive() {
					if err := tx.Where("id = ?", wallet_id).First(&wallet).Update("balance", wallet.Balance.Sub(debitAmount.Amount)).Error; err != nil {
						c.JSON(200, gin.H{
							"error": "Error while processing transaction",
						})
					}
					if err := tx.Where("id = ?", wallet_id).First(&wallet).Update("transaction", Transaction{WalletID: wallet.ID, Type: "debit", Amount: debitAmount.Amount}).Error; err != nil {
						c.JSON(200, gin.H{
							"error": "Error while updating transaction",
						})
					}
					c.JSON(200, gin.H{
						"wallet_balance": wallet.Balance,
					})
				} else {
					// return status 400 if amount to be credited is negative.
					c.JSON(400, gin.H{
						"error": "Cannot use negative value " + debitAmount.Amount.String() + " for operation",
					})
				}
			}
		}
		return nil
	})
}
