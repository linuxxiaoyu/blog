package comment

import (
	"html"
	"net/http"
	"strconv"
	"time"

	"github.com/linuxxiaoyu/blog/pkg/setting"

	"github.com/astaxie/beego/validation"

	"github.com/gin-gonic/gin"
)

// New a comment by article_id and user_id
// POST /comments
// form: aid content
func New(c *gin.Context) {
	aid, _ := strconv.ParseUint(c.PostForm("aid"), 10, 32)
	uid, _ := c.Get("uid")
	content := html.EscapeString(c.PostForm("content"))

	valid := validation.Validation{}
	valid.Required(content, "content").Message("content can't empty")
	if valid.HasErrors() || aid == 0 {
		c.JSON(http.StatusForbidden, nil)
		return
	}

	comment := Comment{
		UserID:    uid.(uint),
		ArticleID: uint(aid),
		Content:   content,
		Time:      time.Now(),
	}

	db := setting.DB
	result := db.Create(&comment)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusForbidden, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": comment.ID,
	})
}
