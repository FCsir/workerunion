package handlers

import (
	"workerunion/db"
	"workerunion/db/models"
)

func CreateSubComment(subComment models.SubComment) error {
	result := db.SqlDB.Create(&subComment)
	return result.Error
}

func FindSubComments(query map[string]interface{}) []models.SubComment {
	var comments []models.SubComment
	queryset := db.SqlDB.Where(query)
	queryset.Find(&comments)

	return comments
}
