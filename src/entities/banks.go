package entities

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

type Banks struct {
	Status  bool       `json:"status"`
	Message string     `json:"message"`
	Data    []BankData `json:"data"`
}

type BankData struct {
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	Code          string    `json:"code"`
	Longcode      string    `json:"longcode"`
	Gateway       string    `json:"gateway"`
	Pay_with_bank bool      `json:"pay_with_bank"`
	Active        bool      `json:"active"`
	Is_deleted    bool      `json:"is_deleted"`
	Country       string    `json:"country"`
	Currency      string    `json:"currency"`
	Type          string    `json:"type"`
	Id            uint      `json:"id"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func GetBankCodeByBankName(bankName string) (code string) {
	data, err := ioutil.ReadFile("/home/salem/Desktop/Franka/src/entities/bnk.json")
	if err != nil {
		fmt.Println(err)
	}

	var banks Banks

	json.Unmarshal([]byte(data), &banks)
	for i, _ := range banks.Data {
		if banks.Data[i].Name == bankName {
			log.Println(banks.Data[i].Code)
			return banks.Data[i].Code
		}
	}
	return "007"
}
