package comment

import (
	"net/http"
	"strconv"

	"github.com/linuxxiaoyu/blog/pkg/setting"

	"github.com/gin-gonic/gin"
)

// Delete a comment by comment_id
// DELETE /comments/:id
// form: token
func Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	uid, _ := c.Get("uid")

	db := setting.DB
	result := db.Where("user_id = ?", uid.(uint)).Delete(&Comment{}, uint(id))
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotImplemented, nil)
		return
	}

	c.JSON(http.StatusOK, nil)
}
