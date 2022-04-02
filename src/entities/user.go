package entities

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string  `json:"username"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	ProfileId string  `json:"profile_id"`
	Profile   Profile `gorm:"foreignKey:ProfileRefer"`
	WalletId  string  `json:"wallet_id"`
	Wallet    Wallet  `gorm:"foreignKey:WalletRefer"`
	Active    bool    `json:"active"`
}

//Generates a password hash for a player's password as storing raw password to db is not ideal
func (u *User) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 18)
	return string(bytes), err
}

// used during login to compare player's login password with the equivalent hash stored in db
func (u *User) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (u *User) activateWallet() (err error) {
	if u.Profile.Address != "" && u.Profile.State != "" {
		u.Wallet.Active = true
		return nil
	}
	return errors.New("Can't activate wallet, atleast two Kyc credentials must be provided")
}
