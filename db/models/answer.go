package models

import "database/sql/driver"

type answerStatus string

const (
	answerPublish answerStatus = "publish"
	answerCancel  answerStatus = "cancel"
)

func (u *answerStatus) Scan(value interface{}) error {
	*u = answerStatus(value.([]byte))
	return nil
}

func (u answerStatus) Value() (driver.Value, error) {
	return string(u), nil
}

type Answer struct {
	BaseModel
	Status answerStatus `gorm:"type:varchar(100)"`

	Content string `gorm:"type:text"`

	UserID uint
	PostID uint

	SubComments []SubComment `gorm:"foreignKey:AnswerID"`
}
