package apis

import (
	"fmt"
	"io"
	"os"
	"time"
	"workerunion/apis/middleware"
	"workerunion/apis/worker_union"

	"github.com/gin-gonic/gin"
)

var Router *gin.Engine = gin.Default()

func init() {

	Router.Use(middleware.Cors())

	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)
	Router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		fmt.Println("----- params:", param.ErrorMessage)
		return fmt.Sprintf("%s  1- [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	Router.Use(gin.Recovery())

	Router.GET("/", worker_union.PopularPosts)
	home := Router.Group("/home")
	{
		home.GET("/recommend", worker_union.RecommendationPosts)
		home.GET("/latest", worker_union.LatestPosts)
		home.GET("/popular", worker_union.PopularPosts)
	}

	auth := Router.Group("/auth")
	{
		auth.POST("/login", worker_union.Login)
		auth.POST("/register", worker_union.Register)
		auth.GET("/activate/:code", worker_union.Activate)
	}

	post := Router.Group("/post")
	{
		post.POST("/:post_id/read", worker_union.ReadPost)
		post.POST("/:post_id/view", worker_union.ViewPost)
		post.Use(middleware.Jwt()).POST("/add", worker_union.AddPost)
		post.POST("/edit", worker_union.EditPost).Use(middleware.Jwt())
		post.POST("/:post_id/detail", worker_union.GetPostDetail).Use(middleware.Jwt())
	}

	user := Router.Group("/user").Use(middleware.Jwt())
	{
		user.GET("/profile", worker_union.Profile)
		user.POST("/save_profile", worker_union.SaveProfile)
		user.GET("/posts", worker_union.UserPosts)
		user.GET("/answers", worker_union.UserAnswers)
		user.GET("/collect", worker_union.UserCollect)
	}
}
