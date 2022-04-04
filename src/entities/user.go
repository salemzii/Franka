package entities

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Profile  Profile
	Wallet   Wallet

	Active bool `json:"active"`
}

//Generates a password hash for a player's password as storing raw password to db is not ideal
func (u *User) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// used during login to compare player's login password with the equivalent hash stored in db
func (u *User) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// activate user wallet
func (u *User) ActivateWallet() (err error) {
	if u.Profile.Address != "" && u.Profile.State != "" {
		u.Wallet.Active = true
		return nil
	}
	return errors.New("can't activate wallet, atleast 2 out of  3 Kyc credentials must be provided")
}

//Query db for a particular player using the player's username
func GetUserObject(username string) (p *User, err error) {
	var user User

	db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("username = ?", username).First(&user).Error; err != nil {
			return err
		}
		return nil
	})
	return &user, nil
}

// hook to create user's wallet on-signup
func (u *User) AfterCreate(tx *gorm.DB) (err error) {

	w := Wallet{
		UserID:  u.ID,
		Balance: decimal.NewFromFloat(0.00),
	}
	if err := tx.Create(&w).Error; err != nil {
		return err
	}
	fmt.Println(w)
	return nil
}

func CreateUser(c *gin.Context) {
	var user User
	// unmarshal incoming json input to user instance
	c.BindJSON(&user)

	// hash user's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		log.Println(err)
	}
	user.Password = string(hashedPassword)

	db.Transaction(func(tx *gorm.DB) error {
		// create user
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		return nil
	})

	// Return json response after saving player
	c.JSON(200, gin.H{
		"username": user.Username,
		"email":    user.Email,
		"password": user.Password,
	})
}

type UpdateUserStruct struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

func UpdateUser(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
	}
	var userUpdate UpdateUserStruct
	var user User
	c.BindJSON(&userUpdate)
	db.Transaction(func(tx *gorm.DB) error {
		if userUpdate.Email != "" && userUpdate.Username != "" {
			tx.Model(&user).Where("id = ?", userId).Select("username", "email").Updates(User{Username: userUpdate.Username, Email: userUpdate.Email})
		} else if userUpdate.Username != "" {
			tx.Model(&user).Where("id = ?", userId).Select("username").Updates(User{Username: userUpdate.Username})
		} else if userUpdate.Email != "" {
			tx.Model(&user).Where("id = ?", userId).Select("email").Updates(User{Email: userUpdate.Email})
		} else {
			c.JSON(400, gin.H{"error": "User update field cannot be empty!"})
		}
		return nil
	})
	c.JSON(200, gin.H{
		"id":       userId,
		"username": userUpdate.Username,
		"email":    userUpdate.Email,
	})
}

func GetUser(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
	}
	var user User
	var wallet Wallet
	db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&user, userId).Error; err != nil {
			return err
		}
		return nil
	})

	if err := db.Where("user_id = ?", user.ID).First(&wallet).Error; err != nil {
		log.Println(err)
	}

	c.JSON(200, gin.H{
		"id":             user.ID,
		"username":       user.Username,
		"email":          user.Email,
		"wallet_id":      wallet.ID,
		"wallet_balance": wallet.Balance,
	})
}
