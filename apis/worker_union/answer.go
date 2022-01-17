package worker_union

import (
	"net/http"
	"workerunion/db/handlers"
	"workerunion/db/models"

	"github.com/gin-gonic/gin"
)

type AddAnswerForm struct {
	Content string `form:"content" binding:"required"`
	PostID  int    `form:"postID" binding:"required"`
}

func AddAnswer(c *gin.Context) {
	var ansForm AddAnswerForm
	err := c.ShouldBind(&ansForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user := c.MustGet("currentUser").(models.User)
	answer := models.Answer{
		Status:  "publish",
		Content: ansForm.Content,
		UserID:  user.ID,
		PostID:  uint(ansForm.PostID),
	}

	if err := handlers.CreateAnswer(answer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"answer": answer})
}
