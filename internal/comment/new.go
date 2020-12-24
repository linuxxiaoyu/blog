package comment

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strconv"
	"time"

	"github.com/linuxxiaoyu/blog/internal/cache"
	"github.com/linuxxiaoyu/blog/internal/setting"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// New a comment by article_id and user_id
// POST /comments
// form: token aid content
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

	cache.Hset("comments", uint(comment.ID), comment)
	cache.Sadd(
		fmt.Sprintf("article_comments:%d", comment.ArticleID),
		strconv.Itoa(int(comment.ID)),
	)

	c.JSON(http.StatusOK, gin.H{
		"id": comment.ID,
	})
}

func GetArticleComments(aid uint) []Comment {
	key := fmt.Sprintf("article_comments:%d", aid)
	comments := []Comment{}

	commentIDs, err := cache.Smembers(key)
	if err != nil {
		db := setting.DB
		result := db.Where("article_id = ?", aid).Find(&comments)
		if result.Error != nil {
			for _, comment := range comments {
				cache.Hset("comments", uint(comment.ID), comment)
				cache.Sadd(
					fmt.Sprintf("article_comments:%d", comment.ArticleID),
					strconv.Itoa(int(comment.ID)),
				)
			}
		}
		return comments
	}

	for _, commentIDStr := range commentIDs {
		commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
		if err != nil {
			return comments
		}
		str, err := cache.Hget("comments", uint(commentID))
		if err != nil {
			return comments
		}
		var comment Comment
		err = json.Unmarshal([]byte(str), &comment)
		if err != nil {
			return comments
		}
		comments = append(comments, comment)
	}

	return comments
}
