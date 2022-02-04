package worker_union

import (
	"fmt"
	"net/http"
	"time"
	"workerunion/db/handlers"
	"workerunion/db/models"

	"github.com/gin-gonic/gin"
)

func Profile(c *gin.Context) {
	user := c.MustGet("currentUser").(models.User)
	profile := map[string]interface{}{
		"email":    user.Email,
		"gender":   user.Gender,
		"nickname": user.NickName,
		"phone":    user.Phone,
		"id":       user.ID,
		"birth":    user.Birth,
		"status":   user.Status,
	}
	c.JSON(http.StatusOK, gin.H{"data": profile})
}

type userProfileForm struct {
	ID     uint   `form:"id" binding:"required"`
	Phone  string `form:"phone"`
	Birth  int64  `form:"birth"`
	Gender string `form:"gender"`
}

func SaveProfile(c *gin.Context) {
	var form userProfileForm
	err := c.ShouldBind(&form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user := c.MustGet("currentUser").(models.User)
	birth := time.Unix(form.Birth/1000, 0)
	birthRFC := birth.Format(time.RFC3339)
	fmt.Println("birth: ", birth, birthRFC)
	data := map[string]interface{}{
		"birth":  birthRFC,
		"gender": form.Gender,
		"phone":  form.Phone,
	}
	handlers.UpdateUser(user, data)

	c.JSON(http.StatusOK, gin.H{"message": "update success"})
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
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
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
