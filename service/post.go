package service

import (
	"workerunion/db/handlers"
	"workerunion/db/models"
)

type PostResult struct {
	*models.Post
	AuthorNickName string `json:"autorNickName"`
}

func getLatestPosts(query map[string]interface{}, orders []map[string]string) {
	posts := handlers.FindPosts(query, orders, 20, 0)
	var userIds []uint
	for _, post := range posts {
		userIds = append(userIds, post.UserID)
	}
	users := handlers.FindUsersByIds(userIds)
	var result []PostResult
	for _, post := range posts {
	}
}
