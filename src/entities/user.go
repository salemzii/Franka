package entities


import (
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
  )
type User struct {
	g
	Username string `json:"username"`

}
