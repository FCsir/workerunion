package worker_union

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"time"
	in "workerunion/apis/internal"
	"workerunion/db/handlers"
	"workerunion/db/models"
	"workerunion/pkg"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Activate(c *gin.Context) {
	code := c.Param("code")

	userEmail, err := pkg.RedisClient.Get(c, "register:"+code).Result()
	fmt.Println("12code-----", userEmail, err)
	if userEmail == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "链接失效, 请重新发送邮件"})
		return
	}

	// activate user status
	users := handlers.FindUsers(map[string]interface{}{"email": userEmail})
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

	existed := handlers.CheckUserByNickname(form.NickName)
	if existed == true {
		c.JSON(http.StatusBadRequest, gin.H{"message": "该昵称已存在"})
		return
	}

	// check the email duplicate
	existed = handlers.CheckUserByEmail(form.Email)
	if existed == true {
		c.JSON(http.StatusBadRequest, gin.H{"message": "该电邮已存在"})
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

	u := uuid.New()
	code := u.String()
	fmt.Println("email---", form.Email)
	err = pkg.RedisClient.Set(c, "register:"+code, form.Email, 10*time.Minute).Err()
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

	hashPassword := pkg.GetMD5Hash(form.Password)

	users := handlers.FindUsers(map[string]interface{}{"email": form.Email})
	if len(users) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "该用户不存在"})
		return
	}

	if users[0].Password != hashPassword {
		c.JSON(http.StatusBadRequest, gin.H{"message": "密码错误"})
		return
	}

	token, err := pkg.GenerateToken(users[0].ID, form.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}
