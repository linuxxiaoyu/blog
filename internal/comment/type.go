package comment

import (
	"time"

	"github.com/linuxxiaoyu/blog/internal/user"
)

// Comment is a table in db.
type Comment struct {
	ID        uint `gorm:"primary_key"`
	UserID    uint `gorm:"forgignKey"`
	ArticleID uint `gorm:"forgignKey"`
	Content   string
	Time      time.Time

	User user.User
}
