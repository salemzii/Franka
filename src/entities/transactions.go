package entities

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	WalletID uint            `json:"transaction_id"`
	Type     string          `json:"type"`
	Amount   decimal.Decimal `json:"amount"`
}

func AllTransactions(c *gin.Context) {
	var tranx []Transaction
	db.Find(&tranx)

	ginH := make(map[int]interface{})
	for i, v := range tranx {
		ginH[i] = v
		fmt.Printf("index: %d, Value %v: ", i, v)
	}
	c.JSON(200, gin.H{
		"all": ginH,
	})
}

func GetWalletTransactions(c *gin.Context) {
	wallet_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
	}
	var tranx []Transaction

	db.Where("wallet_id = ?", wallet_id).Find(&tranx)
	tranxList := make(map[int]interface{})

	for i, v := range tranx {
		tranxList[i] = v
	}
	c.JSON(200, gin.H{
		"your_transactions": tranxList,
	})
}

func GetWalletTransaction(c *gin.Context) {
	wallet_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
	}
	txId, err := strconv.Atoi(c.Param("txId"))
	if err != nil {
		log.Println(err)
	}
	var transaction Transaction
	if err := db.Where("wallet_id = ? AND id = ?", wallet_id, txId).First(&transaction).Error; err == gorm.ErrRecordNotFound {
		c.JSON(200, gin.H{
			"err": "transaction with id " + strconv.Itoa(txId) + " not found!",
		})
		return
	}
	c.JSON(200, gin.H{
		"id":      transaction.ID,
		"created": transaction.CreatedAt,
		"type":    transaction.Type,
		"amount":  transaction.Amount,
	})
}
