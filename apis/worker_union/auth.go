package worker_union

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	in "workerunion/apis/internal"
	"workerunion/db/handlers"
	"workerunion/db/models"
	"workerunion/pkg"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Activate(c *gin.Context) {
	code := c.Param("code")

	log.Println("code-----", code)
	userId := pkg.RedisClient.HGet(c, code, "user_id")
	log.Println("code-----", userId)
	if userId == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "链接失效, 请重新发送邮件"})
		return
	}

	// activate user status
	users := handlers.FindUsers(map[string]interface{}{"id": userId})
	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "user no existed"})
		return
	}

	handlers.ActivateUser(users[0])
	c.JSON(http.StatusOK, gin.H{"message": "账号已激活"})
}

type registerForm struct {
	Email      string `form:"email" binding:"required"`
	NickName   string `form:"nickname" binding:"required"`
	Password   string `form:"password" binding:"required,max=16,min=6"`
	RePassword string `form:"repassword" binding:"required"`
}

func Register(c *gin.Context) {
	var form registerForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if form.Password != form.RePassword {
		c.JSON(http.StatusBadRequest, gin.H{"message": "password conflict"})
		return
	}

	// check the email duplicate
	existed := handlers.CheckUserByEmail(form.Email)
	if existed == true {
		c.JSON(http.StatusBadRequest, gin.H{"message": "该电邮已存在"})
		return
	}

	existed = handlers.CheckUserByNickname(form.NickName)
	if existed == true {
		c.JSON(http.StatusBadRequest, gin.H{"message": "该昵称已存在"})
		return
	}

	t := template.New("activation_account.html")
	homePath, _ := os.Getwd()

	absPath := path.Join(homePath, "static", "email", "activation_account.html")
	t, err := t.ParseFiles([]string{absPath}...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// save account and send email
	hashPassword := pkg.GetMD5Hash(form.Password)
	user := models.User{Email: form.Email, Password: hashPassword, NickName: form.NickName, Status: "register"}
	err = handlers.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	u := uuid.New()
	code := u.String()
	pkg.RedisClient.HSet(c, code, []string{"user_id", "1"}, 10*60)

	log.Println("url---", in.GetHost(c.Request))
	activeUrl := in.GetHost(c.Request) + "/auth/activate/" + code

	data := struct {
		Link string
	}{
		Link: activeUrl,
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	pkg.SendEmail(tpl.String(), "Activation", []string{form.Email})

	c.JSON(http.StatusOK, gin.H{"message": "register ok"})

}

type loginForm struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var form loginForm
	fmt.Println("-----------------hello world")
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// TODO check the user existed

	token, err := pkg.GenerateToken(12, form.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}
