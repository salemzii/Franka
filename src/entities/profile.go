package entities

import (
	"log"
	"strconv"

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
	NinNumber         string `json:"nin"`
	VerificationPhoto string `json:"verificationphoto"`
}

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
