package main

import (
	"github.com/gin-gonic/gin"
	"github.com/linuxxiaoyu/blog/internal/article"
	"github.com/linuxxiaoyu/blog/internal/comment"
	jwt "github.com/linuxxiaoyu/blog/internal/middleware/auth"
	"github.com/linuxxiaoyu/blog/internal/upload"
	"github.com/linuxxiaoyu/blog/internal/user"
)

func handle(r *gin.Engine) {
	gUser := r.Group("/user")
	{
		gUser.POST("", user.SignUp)
		gUser.GET("", user.Login)
	}

	gArticle := r.Group("/articles")
	{
		gArticle.GET("", article.Articles)
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
		gComment.PUT("/:id", comment.Update)
	}

	gUpload := r.Group("/upload")
	gUpload.Use(jwt.Auth)
	{
		gUpload.POST("", upload.Upload)
	}
}
