package entities

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

func init() {
	db, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Wallet{})
	db.AutoMigrate(&Profile{})
	db.AutoMigrate(&Transaction{})
}
