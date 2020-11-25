package article

import (
	"time"

	"github.com/linuxxiaoyu/blog/pkg/comment"
	"github.com/linuxxiaoyu/blog/pkg/user"
)

// Article is a table in db.
type Article struct {
	ID       uint      `gorm:"primary_key"`
	AuthorID uint      `gorm:"forgignKey" json:"author_id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Time     time.Time `json:"time"`

	Author   user.User
	Comments []comment.Comment
}
