package main

import (
	"github.com/gin-gonic/gin"
	"github.com/linuxxiaoyu/blog/internal/middleware/auth"
)

func handle(r *gin.Engine) {
	gUser := r.Group("/user")
	{
		gUser.POST("", register)
		gUser.GET("", login)
	}

	// gArticle := r.Group("/articles")
	// {
	// 	gArticle.GET("", article.Articles)
	// 	gArticle.GET("/:id", article.Get)
	// }
	// gArticle.Use(jwt.Auth)
	// {
	// 	gArticle.POST("", article.New)
	// 	gArticle.DELETE("/:id", article.Delete)
	// 	gArticle.PUT("/:id", article.Update)
	// }

	gComment := r.Group("/comments")
	gComment.Use(auth.Auth)
	{
		gComment.POST("", newComment)
		gComment.DELETE("/:id", deleteComment)
		gComment.PUT("/:id", updateComment)
	}

	// gUpload := r.Group("/upload")
	// gUpload.Use(jwt.Auth)
	// {
	// 	gUpload.POST("", upload.Upload)
	// }
}
