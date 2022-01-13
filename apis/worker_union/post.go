package worker_union

import (
	"log"
	"net/http"
	"workerunion/pkg"

	"github.com/gin-gonic/gin"
)

func ReadPost(c *gin.Context) {
	postId := c.Param("post_id")
	log.Println("post id: ", postId)
	err := pkg.RedisClient.ZIncrBy(c, "post_read_count_sort", 1, "post_id:"+postId).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	count, err := pkg.RedisClient.ZScore(c, "post_read_count_sort", "post_id:"+postId).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	log.Println("-------- post id count: ", count, postId)
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
