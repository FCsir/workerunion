package handlers

import (
	"workerunion/db"
	"workerunion/db/models"
)

func CreateAnswer(answer models.Answer) error {
	result := db.SqlDB.Create(&answer)
	return result.Error
}

func FindAnswers(query map[string]interface{}, orders []map[string]string, limit int, offset int) []models.Answer {
	var answers []models.Answer
	queryset := db.SqlDB.Where(query)
	if limit != 0 {
		queryset = queryset.Limit(limit)
	}
	if offset != 0 {
		queryset = queryset.Offset(offset)
	}
	for _, v := range orders {
		queryset = queryset.Order(v["name"] + " " + v["type"])
	}
	queryset.Find(&answers)

	return answers
}

func CountAnswers(query map[string]interface{}) int64 {
	var result int64
	db.SqlDB.Model(&models.Answer{}).Where(query).Count(&result)
	return result
}
