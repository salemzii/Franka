package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	entity "github.com/salemzii/franka/src/entities"
)

func main() {
	setupServer().Run()
}

// The engine with all endpoints is now extracted from the main function
func setupServer() *gin.Engine {
	//create http router
	router := gin.Default()

	//router.Use(sessions.Sessions("mysession", sessions.NewCookieStore(auth.Secret)))
	router.Use(JSONMiddleware())
	router.GET("/", welcome)

	//Authentication
	router.POST("api/v1/user/register", entity.CreateUser)
	//router.POST("/api/v1/auth/login", auth.LoginFunc)
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
