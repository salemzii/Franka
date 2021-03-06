package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	entity "github.com/salemzii/franka/src/entities"
	"gorm.io/gorm"
)

const userkey = "user"

var Secret = []byte("secret")

type LoginStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthRequired is a simple middleware to check the session cookie
func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		// Abort the request with the appropriate error code
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// Continue down the chain to handler etc
	c.Next()
}

// login is a handler that parses a form and checks for specific data
func LoginFunc(c *gin.Context) {
	session := sessions.Default(c)

	var logindetails LoginStruct

	c.BindJSON(&logindetails)
	if strings.Trim(logindetails.Username, " ") == "" || strings.Trim(logindetails.Password, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Login fields can't be empty"})
		return
	}
	user, err := entity.GetUserObject(logindetails.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{
				"error": "Player with username " + logindetails.Username + " not found",
			})
		}
	}
	fmt.Println(user.Username, user.Password)
	pass, err := user.HashPassword(logindetails.Password)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(pass)
	if logindetails.Username != user.Username || !user.CheckPasswordHash(logindetails.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	// Save the username in the session
	session.Set(userkey, logindetails.Username) // In real world usage you'd set this to the users ID
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated user"})

}

//https://medium.com/wesionary-team/jwt-authentication-in-golang-with-gin-63dbc0816d55#:~:text=JSON%20Web%20Token%20(JWT)%20technology,token%20once%20it%20is%20sent.
