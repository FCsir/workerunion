package worker_union

import (
	"net/http"
	"workerunion/db/handlers"
	"workerunion/db/models"

	"github.com/gin-gonic/gin"
)

type AddSubCommentForm struct {
	Content  string `form:"content" binding:"required"`
	AnswerID int    `form:"answerId" binding:"required"`
}

func AddSubComment(c *gin.Context) {
	var commentForm AddSubCommentForm
	err := c.ShouldBind(&commentForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user := c.MustGet("currentUser").(models.User)
	comment := models.SubComment{
		Status:   "publish",
		Content:  commentForm.Content,
		UserID:   user.ID,
		AnswerID: uint(commentForm.AnswerID),
	}

	if err := handlers.CreateSubComment(comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"comment": comment})
}
