package comment

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/linuxxiaoyu/blog/internal/cache"
	"github.com/linuxxiaoyu/blog/internal/setting"

	"github.com/gin-gonic/gin"
)

// Delete a comment by comment_id
// DELETE /comments/:id
// form: token
func Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	uid, _ := c.Get("uid")

	db := setting.DB
	comment := Comment{}
	result := db.Where("user_id = ?", uid.(uint)).Delete(&comment, uint(id))
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotImplemented, nil)
		return
	}

	commentStr, err := cache.Hget("comments", uint(id))
	if err == nil {
		if json.Unmarshal([]byte(commentStr), &comment) == nil {
			cache.Srem(
				fmt.Sprintf("article_comments:%d", comment.ArticleID),
				strconv.Itoa(int(comment.ID)),
			)
		}
	}
	cache.Hdel("comments", uint(id))

	c.JSON(http.StatusOK, nil)
}
