package worker_union

import (
	"log"
	"net/http"
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

}
