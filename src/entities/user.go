package entities

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Profile  Profile `gorm:"foreignKey:ProfileRefer"`

	Wallet Wallet `gorm:"foreignKey:WalletRefer"`
	Active bool   `json:"active"`
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

// activate user wallet
func (u *User) activateWallet() (err error) {
	if u.Profile.Address != "" && u.Profile.State != "" {
		u.Wallet.Active = true
		return nil
	}
	return errors.New("can't activate wallet, atleast 2 out of  3 Kyc credentials must be provided")
}

// hook to create user's wallet on-signup
func (u *User) AfterCreate(tx *gorm.DB) (err error) {

	w := Wallet{
		Balance: decimal.NewFromFloat(0.00),
	}
	if err := tx.Create(&w).Error; err != nil {
		return err
	}
	return nil
}

func CreateUser(c *gin.Context) {
	var user User

	// unmarshal incoming json input to user instance
	c.BindJSON(&user)

	db.Transaction(func(tx *gorm.DB) error {
		pass, err := user.HashPassword(user.Password)
		if err != nil {
			log.Println(err)
		}
		// set user password to the generated password hash
		user.Password = pass
		// create user
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		// Return json response after saving player
		c.JSON(200, gin.H{
			"username": user.Username,
			"email":    user.Email,
			"password": user.Password,
		})
		return nil
	})
}
