package worker_union

import (
	"log"
	"net/http"
	"strconv"
	"workerunion/db/handlers"
	"workerunion/db/models"
	main_in "workerunion/internal"
	"workerunion/pkg"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {

}

func LatestPosts(c *gin.Context) {
	var query = map[string]interface{}{
		"status": "publish",
	}
	var orders = []map[string]string{
		{
			"type": "desc",
			"name": "created_at",
		},
	}
	posts := handlers.FindPosts(query, orders, 20, 0)
	c.JSON(http.StatusOK, gin.H{"data": posts})
}

func PopularPosts(c *gin.Context) {
	postIds, err := pkg.RedisClient.ZRandMember(c, "post_read_count_sort", 30, true).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	log.Println("post ids: ", postIds)
	var posts = []models.Post{}

	if postIds != nil || len(postIds) != 0 {
		ids, err := main_in.StringArray2IntArray(postIds)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		posts = handlers.FindPostsByIds(ids)
	}
	c.JSON(http.StatusOK, gin.H{"data": posts})
}

func ViewPost(c *gin.Context) {
	postIdStr := c.Param("post_id")
	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	posts := handlers.FindPostsByIds([]int{postId})
	if len(posts) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}
	post := posts[0]
	if post.Status != "publish" {
		c.JSON(http.StatusNotFound, gin.H{"message": "close post"})
		return
	}

	query := map[string]interface{}{
		"post_id": post.ID,
		"status":  "publish",
	}
	orders := []map[string]string{
		{"created_at": "desc"},
	}
	answers := handlers.FindAnswers(query, orders, 0, 0)
	var answerIds []int
	for _, answer := range answers {
		answerIds = append(answerIds, int(answer.ID))
	}

	comments := handlers.FindSubCommentsByAnswerIds(answerIds)

	c.JSON(http.StatusNotFound, gin.H{"post": post, "answers": answers, "comments": comments})

}
