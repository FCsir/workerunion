package models

type Post struct {
	BaseModel
	Status string

	Content   string `gorm:"type:text"`
	Title     string
	Readcount uint

	UserID uint

	Comments []Comment `gorm:"foreignKey:PostID"`
}
