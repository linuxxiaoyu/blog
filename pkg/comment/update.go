package comment

import (
	"fmt"
	"html"
	"net/http"
	"strconv"
	"time"

	"github.com/linuxxiaoyu/blog/pkg/cache"

	"github.com/linuxxiaoyu/blog/pkg/setting"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// Update a comment by comment_id
// PUT /comments/:id
// form: token content
func Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	uid, _ := c.Get("uid")

	content := html.EscapeString(c.PostForm("content"))
	valid := validation.Validation{}
	valid.Required(content, "content").Message("content can't empty")
	if valid.HasErrors() {
		c.JSON(http.StatusForbidden, nil)
		return
	}

	var comment Comment
	db := setting.DB
	result := db.First(&comment, uint(id))
	if result.RowsAffected == 0 || comment.UserID != uid.(uint) {
		c.JSON(http.StatusNotImplemented, nil)
		return
	}

	comment.Content = content
	comment.Time = time.Now().Local()
	db.Save(&comment)

	cache.Hset("comments", uint(id), comment)
	cache.Sadd(
		fmt.Sprintf("article_comments:%d", comment.ArticleID),
		strconv.Itoa(int(id)),
	)

	c.JSON(http.StatusOK, nil)
}
