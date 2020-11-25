package article

import (
	"html"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/linuxxiaoyu/blog/pkg/setting"
)

// Update an article
// PUT /articles/:id
// form: token title content
func Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	uid, _ := c.Get("uid")

	db := setting.DB
	var article Article
	result := db.First(&article, uint(id))
	if result.RowsAffected == 0 || article.AuthorID != uid.(uint) {
		c.JSON(http.StatusForbidden, nil)
		return
	}

	title := html.EscapeString(c.PostForm("title"))
	if title != "" {
		article.Title = title
	}

	content := html.EscapeString(c.PostForm("content"))
	if content != "" {
		article.Content = content
	}

	article.Time = time.Now().Local()

	result = db.Save(&article)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotImplemented, nil)
		return
	}

	c.JSON(http.StatusOK, nil)
}
