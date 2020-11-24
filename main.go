package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/linuxxiaoyu/blog/pkg/setting"
)

func main() {
	setting.Init()
	if setting.IsReleaseMode() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		// time.Sleep(10 * time.Second)
		c.String(http.StatusOK, "ok")
	})
	graceful(r)
}
