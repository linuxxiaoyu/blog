package article

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/linuxxiaoyu/blog/pkg/cache"
	"github.com/linuxxiaoyu/blog/pkg/setting"
)

// Articles return 10 articles order by id
// GET /articles
func Articles(c *gin.Context) {
	db := setting.DB

	articles := []Article{}

	result := db.Preload("Author").Preload("Comments.User").Order("time desc").Find(&articles).Limit(10)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	datas := []gin.H{}
	for _, article := range articles {
		datas = append(datas, articleResponse(&article))
		cache.Hset("articles", article.ID, article)
	}

	c.JSON(http.StatusOK, datas)
}
