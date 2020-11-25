package article

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/linuxxiaoyu/blog/pkg/setting"
)

// Get an article
// GET /articles/:id
func Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	article := Article{
		ID: uint(id),
	}

	db := setting.DB
	result := db.Preload("Author").Preload("Comments.User").First(&article)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, articleResponse(&article))
}

func articleResponse(article *Article) gin.H {
	if article == nil {
		return gin.H{}
	}

	layout := "2006-01-02 15:04:05"

	comments := []gin.H{}
	for _, comment := range article.Comments {
		comments = append(comments, gin.H{
			"id":       comment.ID,
			"username": comment.User.Name,
			"content":  comment.Content,
			"time":     comment.Time.Format(layout),
		})
	}

	return gin.H{
		"id":       article.ID,
		"title":    article.Title,
		"author":   article.Author.Name,
		"content":  article.Content,
		"time":     article.Time.Format(layout),
		"comments": comments,
	}
}
