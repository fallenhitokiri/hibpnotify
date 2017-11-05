package hibpnotify

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Breach struct {
	gorm.Model

	Title        string
	Name         string
	Domain       string
	AddedDate    time.Time
	ModifiedDate time.Time
	Description  string
	Users        []User `gorm:"many2many:user_breaches;"`
}
