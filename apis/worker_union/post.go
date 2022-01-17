package worker_union

import (
	"log"
	"net/http"
	"strconv"
	"workerunion/db/handlers"
	"workerunion/db/models"
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

type AddPostForm struct {
	Content string `form:"content"`
	Title   string `form:"title" binding:"required"`
}

func AddPost(c *gin.Context) {
	var postForm AddPostForm
	err := c.ShouldBind(&postForm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	user := c.MustGet("currentUser").(models.User)
	post := models.Post{
		Status:    "draft",
		Content:   postForm.Content,
		Title:     postForm.Title,
		Readcount: 0,
		UserID:    user.ID,
	}

	if err := handlers.CreatePost(post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"post": post})
}

type EditPostForm struct {
	Content string `form:"content"`
	Title   string `form:"title" binding:"required"`
	ID      int    `form:"id" binding:"required"`
	Status  string `form:"status" binding:"required"`
}

func EditPost(c *gin.Context) {
	var postForm EditPostForm
	err := c.ShouldBind(&postForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// user := c.MustGet("currentUser").(models.User)
	posts := handlers.FindPostsByIds([]int{postForm.ID})
	if len(posts) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	post := posts[0]

	user := c.MustGet("currentUser").(models.User)
	if post.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"message": "forbidden"})
		return
	}

	data := map[string]interface{}{
		"content": postForm.Content,
		"title":   postForm.Title,
		"status":  postForm.Status,
	}

	handlers.UpdatePost(post, data)
	c.JSON(http.StatusOK, gin.H{"post": post})
}

func GetPostDetail(c *gin.Context) {
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
	user := c.MustGet("currentUser").(models.User)
	if post.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"message": "forbidden"})
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"post": post})
}
