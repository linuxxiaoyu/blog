package article

import (
	"html"
	"net/http"
	"time"

	"github.com/linuxxiaoyu/blog/pkg/setting"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// New create a article.
// POST /article
// form: token title content
func New(c *gin.Context) {
	title := html.EscapeString(c.PostForm("title"))
	content := html.EscapeString(c.PostForm("content"))

	valid := validation.Validation{}
	valid.Required(title, "title").Message("title can't empty")
	valid.Required(content, "content").Message("content can't empty")

	if valid.HasErrors() {
		c.JSON(http.StatusForbidden, nil)
		return
	}

	uid, _ := c.Get("uid")
	article := Article{
		AuthorID: uid.(uint),
		Title:    title,
		Content:  content,
		Time:     time.Now().Local(),
	}

	db := setting.DB
	result := db.Create(&article)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": article.ID,
	})
}
