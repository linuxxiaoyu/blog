package jwt

import (
	"net/http"
	"time"

	"github.com/linuxxiaoyu/blog/internal/auth"
	"github.com/linuxxiaoyu/blog/internal/cache"
	"github.com/linuxxiaoyu/blog/internal/setting"
	"github.com/linuxxiaoyu/blog/internal/user"

	"github.com/gin-gonic/gin"
)

func forbidden(c *gin.Context) {
	c.JSON(http.StatusForbidden, nil)
	c.Abort()
}

// form: token
func Auth(c *gin.Context) {
	token := c.PostForm("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, nil)
		c.Abort()
		return
	}

	claims, err := auth.ParseToken(token)
	if err != nil || claims.ExpiresAt < time.Now().Unix() {
		cache.Srem("tokens", token)
		forbidden(c)
		return
	}

	c.Set("uid", claims.ID)

	var user user.User
	db := setting.DB
	result := db.First(&user, claims.ID)
	if result.RowsAffected == 0 || user.Name != claims.Name {
		forbidden(c)
		return
	}

	cache.Sadd("tokens", token)

	c.Next()
}
