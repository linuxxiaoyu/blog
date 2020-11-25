package main

import (
	"github.com/gin-gonic/gin"
	"github.com/linuxxiaoyu/blog/pkg/article"
	"github.com/linuxxiaoyu/blog/pkg/comment"
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
	{
		gArticle.GET("/", article.Articles)
		gArticle.GET("/:id", article.Get)
	}
	gArticle.Use(jwt.Auth)
	{
		gArticle.POST("", article.New)
		gArticle.DELETE("/:id", article.Delete)
		gArticle.PUT("/:id", article.Update)
	}

	gComment := r.Group("/comments")
	gComment.Use(jwt.Auth)
	{
		gComment.POST("", comment.New)
		gComment.DELETE("/:id", comment.Delete)
	}
}
