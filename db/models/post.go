package models

import "database/sql/driver"

type postStatus string

const (
	publish postStatus = "publish"
	draft   postStatus = "draft"
	cancel  postStatus = "cancel"
)

func (u *postStatus) Scan(value interface{}) error {
	*u = postStatus(value.([]byte))
	return nil
}

func (u postStatus) Value() (driver.Value, error) {
	return string(u), nil
}

type Post struct {
	BaseModel
	Status postStatus `sql:"type:post_status"`

	Content   string `gorm:"type:text"`
	Title     string
	Readcount uint

	UserID uint

	Comments []Comment `gorm:"foreignKey:PostID"`
}
