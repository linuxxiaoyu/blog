package article

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/linuxxiaoyu/blog/internal/cache"
	"github.com/linuxxiaoyu/blog/internal/setting"
)

// Delete an article
// DELETE /articles/:id
// form: token
func Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	uid, _ := c.Get("uid")

	db := setting.DB
	result := db.Where("author_id = ?", uid.(uint)).Delete(&Article{}, uint(id))
	if result.RowsAffected == 0 {
		c.JSON(http.StatusForbidden, nil)
		return
	}

	cache.Hdel("articles", uint(id))

	c.JSON(http.StatusOK, nil)
}
