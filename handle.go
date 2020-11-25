package main

import (
	"github.com/gin-gonic/gin"
	"github.com/linuxxiaoyu/blog/pkg/user"
)

func handle(r *gin.Engine) {
	gUser := r.Group("/user")
	{
		gUser.POST("", user.SignUp)
	}
}
