package main

import (
	"github.com/gin-gonic/gin"
	"github.com/linuxxiaoyu/blog/pkg/article"
	jwt "github.com/linuxxiaoyu/blog/pkg/middleware/auth"
	"github.com/linuxxiaoyu/blog/pkg/user"
)

func handle(r *gin.Engine) {
	gUser := r.Group("/user")
	{
		gUser.POST("", user.SignUp)
		gUser.GET("", user.Login)
	}

	gArticle := r.Group("/articles")
	gArticle.Use(jwt.Auth)
	{
		gArticle.POST("", article.New)
		gArticle.DELETE("/:id", article.Delete)
		gArticle.PUT("/:id", article.Update)
	}
}
