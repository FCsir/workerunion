package models

type User struct {
	BaseModel
	NickName string
	Email    string
	Password string
	// register active
	Status string

	Posts       []Post       `gorm:"foreignKey:UserID"`
	Comments    []Comment    `gorm:"foreignKey:UserID"`
	SubComments []SubComment `gorm:"foreignKey:UserID"`
}
