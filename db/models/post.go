package models

import "database/sql/driver"

type PostStatus string

const (
	publish PostStatus = "publish"
	draft   PostStatus = "draft"
	cancel  PostStatus = "cancel"
)

func (u *PostStatus) Scan(value interface{}) error {
	*u = PostStatus(value.([]byte))
	return nil
}

func (u PostStatus) Value() (driver.Value, error) {
	return string(u), nil
}

type Post struct {
	BaseModel
	Status PostStatus `gorm:"type:varchar(100)" json:"status"`

	Content   string `gorm:"type:text" json:"content"`
	Title     string `gorm:"type:varchar(2000)" json:"title"`
	Readcount uint   `json:"readCount"`
	IsAnymous bool
	Tags      string `gorm:"type:varchar(2000)" json:"tags"`

	UserID uint `json:"userID"`

	Answers []Answer `gorm:"foreignKey:PostID"`
}
