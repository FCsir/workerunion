package worker_union

import (
	"log"
	"net/http"
	"workerunion/pkg"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {

}

func LatestPosts(c *gin.Context) {
	type post struct {
		Content string
		Title   string
	}
	var posts = []post{
		{Content: "hello", Title: "ceshi"},
	}
	c.JSON(http.StatusOK, gin.H{"data": posts})
}

func PopularPosts(c *gin.Context) {
	postIds, err := pkg.RedisClient.ZRandMember(c, "post_read_count_sort", 30, true).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// get posts info
	log.Println("post ids: ", postIds)
	c.JSON(http.StatusOK, gin.H{"data": postIds})
}
