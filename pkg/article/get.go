package article

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/linuxxiaoyu/blog/pkg/cache"

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

	var ginH gin.H
	// FIXME 获取到articles后，是否需要获取comments ？
	articleStr, err := cache.Hget("articles", uint(id))
	if err == nil && json.Unmarshal([]byte(articleStr), &article) == nil {
		ginH = articleResponse(&article)
		c.JSON(http.StatusOK, ginH)
		return
	}

	db := setting.DB
	result := db.Preload("Author").Preload("Comments.User").First(&article)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	cache.Hset("articles", uint(id), article)

	ginH = articleResponse(&article)
	c.JSON(http.StatusOK, ginH)
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

		cache.Hset("comments", comment.ID, comment)
		cache.Sadd(
			fmt.Sprintf("article_comments:%d", comment.ArticleID),
			strconv.Itoa(int(comment.ID)),
		)
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
