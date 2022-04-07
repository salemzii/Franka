package entities

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Profile struct {
	gorm.Model
	UserID        uint
	PhoneNumber   string `json:"phone_number"`
	Address       string `json:"home_address"`
	BankName      string `json:"bank_name"`
	AccountNumber string `json:"account_number"`
	BvnNumber     string `json:"bvn"`
	// VerificationPhoto string `json:"verificationphoto"`
}

func (profile *Profile) HasPhaseOneKyc() bool {
	if profile.PhoneNumber != "" && profile.Address != "" {
		return true
	}
	return false
}
func (profile *Profile) HasAllKyc() bool {
	if profile.HasPhaseOneKyc() && profile.AccountNumber != "" && profile.BvnNumber != "" {
		return true
	}
	return false
}

func (profile *Profile) AfterCreate(tx *gorm.DB) error {
	var user User
	if err := tx.Model(&user).Where("id = ?", profile.UserID).Error; err == gorm.ErrRecordNotFound {
		return gorm.ErrRecordNotFound
	}
	url := "https://api.paystack.co/customer/" + "CUS_mkkqn2euu9d8jul" + "/identification"
	fmt.Println(user.Customer_code)
	bank_code := GetBankCodeByBankName(profile.BankName)
	if profile.HasAllKyc() {
		// make a request to paystack to validate user credentials
		ValidateCustomerCred(url, profile.AccountNumber, profile.BvnNumber, bank_code, user.Username)
	}
	return nil
}

func ValidateCustomerCred(url string, account_number string, bvn string, bank_code string, username string) {

	type CustomerCredentials struct {
		Country        string `json:"string"`
		Type           string `json:"type"`
		Account_number string `json:"account_number"`
		Bvn            string `json:"bvn"`
		Bank_code      string `json:"bank_code"`
		First_name     string `json:"first_name"`
		Last_name      string `json:"last_name"`
	}
	body := CustomerCredentials{
		Country:        "NG",
		Type:           "bank_account",
		Account_number: account_number,
		Bvn:            bvn,
		Bank_code:      bank_code, // provide an api that fetches bank codes
		First_name:     username,
		Last_name:      "franka",
	}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)

	client := http.Client{}

	req, err := http.NewRequest("POST", url, payloadBuf)
	if err != nil {
		log.Fatal(err)
	}

	req.Header = http.Header{
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{"BEARER sk_live_346541ab25e7ff4b0348aa1e1e06bfcd866e2ce2"},
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res.StatusCode, res.Status)
	defer res.Body.Close()
}

//https://paystack.com/docs/identity-verification/validate-customer/#listen-for-verification-status
type Customer struct {
	//gorm.Model
	//UserID  uint
	Status  bool       `json:"status"`
	Message string     `json:"message"`
	Data    DataStruct `json:"data"`
}

type DataStruct struct {
	//gorm.Model
	//CustomerID      uint
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
	kyc.UserID = uint(userId)

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
	c.JSON(200, gin.H{
		"username":   user.Username,
		"email":      user.Email,
		"address":    kyc.Address,
		"accountNum": kyc.AccountNumber,
		"phoneNum":   kyc.PhoneNumber,
		"createdAt":  kyc.CreatedAt,
	})
}
