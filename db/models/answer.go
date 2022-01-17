package models

type Answer struct {
	BaseModel
	Status string

	Content string `gorm:"type:text"`

	UserID uint
	PostID uint

	SubComments []SubComment `gorm:"foreignKey:AnswerID"`
}
