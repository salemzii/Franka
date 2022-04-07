package main

import (
	"fmt"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	auth "github.com/salemzii/franka/src/auth"
	entity "github.com/salemzii/franka/src/entities"
)

func main() {

	setupServer().Run()
}

// The engine with all endpoints is now extracted from the main function
func setupServer() *gin.Engine {

	//create http router
	router := gin.Default()

	router.Use(sessions.Sessions("mysession", sessions.NewCookieStore(auth.Secret)))
	router.Use(JSONMiddleware())
	router.GET("/", welcome)

	//Authentication
	router.POST("api/v1/auth/register", entity.CreateUser)
	router.POST("/api/v1/auth/login", auth.LoginFunc)

	router.PUT("api/v1/user/:id/update", entity.UpdateUser)
	router.GET("api/v1/user/:id", entity.GetUser)
	router.POST("api/v1/wallet/:id/credit", entity.CreditWallet)
	router.POST("api/v1/wallet/:id/debit", entity.DebitWallet)
	router.GET("api/v1/wallet/:id", entity.GetWallet)
	router.GET("api/v1/transactions", entity.AllTransactions)
	router.GET("api/v1/:id/transactions", entity.GetWalletTransactions)
	router.GET("api/v1/:id/transaction/:txId", entity.GetWalletTransaction)
	router.POST("api/v1/:id/kyc", entity.AddKyc)
	//router.GET("/api/v1/auth/logout", auth.Logout)

	// Private group, require authentication to access any wallet resources
	private := router.Group("/private")
	// private.Use(auth.AuthRequired)
	{
	}

	fmt.Println(private)

	return router
}

func JSONMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Accept", "application/json")
		c.Next()
	}
}

func welcome(c *gin.Context) {

	c.JSON(200, gin.H{
		"message": "Hello welcome to Franka api",
	})
}
