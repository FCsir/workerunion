package models

import (
	"database/sql/driver"
	"time"
)

type userStatus string

const (
	active   userStatus = "active"
	register userStatus = "register"
)

func (u *userStatus) Scan(value interface{}) error {
	*u = userStatus(value.([]byte))
	return nil
}

func (u userStatus) Value() (driver.Value, error) {
	return string(u), nil
}

type User struct {
	BaseModel
	NickName string     `gorm:"type:varchar(100)"`
	Email    string     `gorm:"type:varchar(100)"`
	Password string     `gorm:"type:varchar(3000)"`
	Gender   string     `gorm:"type:varchar(100)"`
	Phone    string     `gorm:"type:varchar(100)"`
	Birth    *time.Time `json:"birth"`
	// register active
	Status userStatus `gorm:"type:varchar(100)"`

	Posts       []Post       `gorm:"foreignKey:UserID"`
	Answers     []Answer     `gorm:"foreignKey:UserID"`
	SubComments []SubComment `gorm:"foreignKey:UserID"`
}
