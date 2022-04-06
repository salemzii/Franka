package entities

import (
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Profile struct {
	gorm.Model
	UserID            uint
	PhoneNumber       string `json:"phone_number"`
	Address           string `json:"home_address"`
	State             string `json:"state"`
	AccountNumber     string `json:"account_number"`
	BvnNumber         string `json:"bvn"`
	VerificationPhoto string `json:"verificationphoto"`
}

type Customer struct {
	gorm.Model
	UserID  uint
	Status  bool       `json:"status"`
	Message string     `json:"message"`
	Data    DataStruct `json:"data"`
}

type DataStruct struct {
	gorm.Model
	CustomerID      uint
	Email           string    `json:"email"`
	Integration     uint      `json:"integration"`
	Domain          string    `json:"domain"`
	Customer_code   string    `json:"customer_code"`
	Id              uint      `json:"id"`
	Identified      bool      `json:"identified"`
	Identifications string    `json:"identifications"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// sk_live_346541ab25e7ff4b0348aa1e1e06bfcd866e2ce2

func AddKyc(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
	}
	var user User
	var kyc Profile

	c.BindJSON(&kyc)

	db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&user, userId).Error; err != nil {
			return err
		}
		if err := tx.Create(&kyc).Error; err != nil {
			log.Println(err)
			return err
		}
		return nil
	})
}
