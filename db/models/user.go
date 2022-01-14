package models

import "database/sql/driver"

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
	NickName string
	Email    string
	Password string
	// register active
	Status userStatus `sql:"type:user_status"`

	Posts       []Post       `gorm:"foreignKey:UserID"`
	Comments    []Comment    `gorm:"foreignKey:UserID"`
	SubComments []SubComment `gorm:"foreignKey:UserID"`
}
