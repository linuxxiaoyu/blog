package main

import (
	"github.com/gin-gonic/gin"
	"github.com/linuxxiaoyu/blog/pkg/setting"
)

func main() {
	setting.Init()

	r := gin.Default()
	r.Run(":8080")
}
