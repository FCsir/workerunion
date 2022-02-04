package models

import "database/sql/driver"

type subCommentStatus string

const (
	subCommentPublish subCommentStatus = "publish"
	subCommentCancel  subCommentStatus = "cancel"
)

func (u *subCommentStatus) Scan(value interface{}) error {
	*u = subCommentStatus(value.([]byte))
	return nil
}

func (u subCommentStatus) Value() (driver.Value, error) {
	return string(u), nil
}

type SubComment struct {
	BaseModel
	Status subCommentStatus `gorm:"type:varchar(100)"`

	Content string `gorm:"type:text"`

	UserID   uint
	AnswerID uint
}
