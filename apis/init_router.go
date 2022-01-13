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

	Router.Use(gin.Logger())
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)
	Router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
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

	// TODO add recovery

	Router.GET("/", worker_union.PopularPosts)
	home := Router.Group("/home")
	{
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
	}

	user := Router.Group("/user").Use(middleware.Jwt())
	{
		user.GET("/profile", worker_union.Profile)
		user.GET("/posts", worker_union.UserPosts)
		user.GET("/comments", worker_union.UserComments)
		user.GET("/collect", worker_union.UserCollect)
	}
}
