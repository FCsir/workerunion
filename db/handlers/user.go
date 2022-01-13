package handlers

import (
	"workerunion/db"
	"workerunion/db/models"
)

func CheckUserByEmail(email string) bool {
	var users []models.User
	db.SqlDB.Where("email = ?", email).Find(&users)
	if len(users) == 0 {
		return false
	}
	return true
}

func CheckUserByNickname(nickname string) bool {
	var users []models.User
	db.SqlDB.Where("nickname = ?", nickname).Find(&users)
	if len(users) == 0 {
		return false
	}
	return true
}

func CreateUser(user models.User) error {
	result := db.SqlDB.Create(&user)

	return result.Error
}

func FindUsers(userQuery map[string]interface{}) []models.User {
	var users []models.User
	db.SqlDB.Where(userQuery).Find(&users)
	return users
}

func ActivateUser(user models.User) {
	db.SqlDB.Model(&user).Update("status", "active")
}
