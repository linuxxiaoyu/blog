package data

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/linuxxiaoyu/blog/internal/cache"
	"github.com/linuxxiaoyu/blog/internal/setting"
)

// Comment is a table in db
type Comment struct {
	ID      uint32 `gorm:"primary_key"`
	UID     uint32 `gorm:"column:user_id"`
	AID     uint32 `gorm:"column:article_id"`
	Content string
	Time    time.Time
}

func CreateComment(ctx context.Context, comment *Comment) (uint32, error) {
	if comment == nil {
		return 0, errors.New("comment is nil")
	}
	if comment.UID <= 0 || comment.AID <= 0 {
		return 0, errors.New("uid and aid can't smaller than 0")
	}
	db := setting.DB().WithContext(ctx)
	result := db.Create(&comment)
	if result.Error != nil {
		return 0, errors.WithStack(result.Error)
	}

	cache.Hset(ctx, "comments", comment.ID, comment)
	cache.Sadd(
		ctx,
		fmt.Sprintf("article_comments:%d", comment.AID),
		strconv.Itoa(int(comment.ID)),
	)
	return comment.ID, nil
}

func DeleteComment(ctx context.Context, id, uid uint32) error {
	if id <= 0 || uid <= 0 {
		return errors.New("uid and id can't smaller than 0")
	}

	tx := setting.DB().WithContext(ctx)
	comment := Comment{}
	result := tx.Where("user_id = ?", uid).Delete(&comment, uint(id))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected <= 0 {
		return errors.New("not found")
	}

	commentStr, err := cache.Hget(ctx, "comments", id)
	if err == nil {
		if json.Unmarshal([]byte(commentStr), &comment) == nil {
			cache.Srem(
				fmt.Sprintf("article_comments:%d", comment.AID),
				strconv.Itoa(int(comment.ID)),
			)
		}
	}
	cache.Hdel(ctx, "comments", id)
	return nil
}

func UpdateComment(ctx context.Context, comment *Comment) error {
	if comment == nil {
		return errors.New("comment is nil")
	}
	if comment.UID <= 0 || comment.ID <= 0 {
		return errors.New("uid and id can't smaller than 0")
	}

	tx := setting.DB().WithContext(ctx)
	result := tx.UpdateColumns(&comment)
	if result.Error != nil {
		return result.Error
	}

	cache.Hset(ctx, "comments", comment.ID, comment)
	cache.Sadd(
		ctx,
		fmt.Sprintf("article_comments:%d", comment.AID),
		strconv.Itoa(int(comment.ID)),
	)
	return nil
}
