package upload

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Upload a file
// POST /upload
// form: token file
func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusForbidden, nil)
		return
	}

	c.SaveUploadedFile(file, "files/"+file.Filename)
	c.JSON(http.StatusOK, nil)
}
