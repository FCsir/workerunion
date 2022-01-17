package models

type SubComment struct {
	BaseModel
	Status string

	Content string `gorm:"type:text"`

	UserID   uint
	AnswerID uint
}
