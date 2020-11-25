package article

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func New(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}
