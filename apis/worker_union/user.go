package worker_union

import (
	"net/http"
	"workerunion/db/handlers"
	"workerunion/db/models"

	"github.com/gin-gonic/gin"
)

func Profile(c *gin.Context) {
	user := c.MustGet("currentUser").(models.User)
	c.JSON(http.StatusOK, gin.H{"data": user})
}

type userPostQuery struct {
	PageSize int `form:"pageSize" binding:"required"`
	PageNum  int `form:"pageNum" binding:"required"`
}

func UserPosts(c *gin.Context) {
	var postQuery userPostQuery

	user := c.MustGet("currentUser").(models.User)
	err := c.ShouldBindQuery(&postQuery)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": err.Error()})
		return
	}

	query := map[string]interface{}{
		"user_id": user.ID,
	}
	orders := []map[string]string{
		{"created_at": "desc"},
	}
	posts := handlers.FindPosts(query, orders, postQuery.PageSize, postQuery.PageNum)

	c.JSON(http.StatusOK, gin.H{"posts": posts, "total": handlers.CountPosts(query)})
}

type UserAnswerQuery struct {
	PageSize int `form:"pageSize" binding:"required"`
	PageNum  int `form:"pageNum" binding:"required"`
}

func UserAnswers(c *gin.Context) {
	var answersQuery UserAnswerQuery

	user := c.MustGet("currentUser").(models.User)
	err := c.ShouldBindQuery(&answersQuery)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": err.Error()})
		return
	}

	query := map[string]interface{}{
		"user_id": user.ID,
	}
	orders := []map[string]string{
		{"created_at": "desc"},
	}
	answers := handlers.FindAnswers(query, orders, answersQuery.PageSize, answersQuery.PageNum)

	c.JSON(http.StatusOK, gin.H{"posts": answers, "total": handlers.CountAnswers(query)})

}

func UserCollect(c *gin.Context) {

}
