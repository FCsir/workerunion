package handlers

import (
	"workerunion/db"
	"workerunion/db/models"
)

func FindPosts(query map[string]interface{}, orders []map[string]string, limit int, offset int) []models.Post {
	var posts []models.Post
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
	queryset.Find(&posts)

	return posts
}

func FindPostsByIds(ids []int) []models.Post {
	var posts []models.Post
	db.SqlDB.Find(&posts, ids)

	return posts
}
