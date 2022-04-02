package entities

import "gorm.io/gorm"

type Profile struct {
	gorm.Model
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"home_address`
	State       string `json:"state"`
	ValidIdUrl  string `json:"valididurl"`
}
