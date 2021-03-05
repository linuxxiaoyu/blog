package main

import (
	"github.com/gin-gonic/gin"
	"github.com/linuxxiaoyu/blog/internal/setting"
)

func main() {
	setting.Init()
	if setting.IsReleaseMode() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	handle(r)

	graceful(r)
}
