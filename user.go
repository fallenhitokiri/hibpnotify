package hibpnotify

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model

	Email    string
	Breaches []Breach `gorm:"many2many:user_breaches"`
}
