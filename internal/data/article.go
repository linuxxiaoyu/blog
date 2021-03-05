package data

import (
	"context"
	"errors"
	"time"

	"github.com/linuxxiaoyu/blog/internal/cache"
	"github.com/linuxxiaoyu/blog/internal/setting"
)

type Article struct {
	ID       uint32 `gorm:"primary_key"`
	AuthorID uint32
	Title    string
	Content  string
	Time     time.Time
}

func GetArticles(ctx context.Context) ([]Article, error) {
	articles := []Article{}
	tx := setting.DB().WithContext(ctx)
	result := tx.Order("time desc").Find(&articles).Limit(10)
	if result.Error != nil {
		return nil, result.Error
	}

	return articles, nil
}

func GetArticle(ctx context.Context, id uint32) (Article, error) {
	article := Article{ID: id}
	tx := setting.DB().WithContext(ctx)
	result := tx.First(&article)
	if result.Error != nil {
		return article, result.Error
	}
	if result.RowsAffected <= 0 {
		return article, errors.New("not found")
	}

	cache.Hset(ctx, "articles", id, article)
	return article, nil
}

func CreateArticle(ctx context.Context, uid uint32, title, content string, time time.Time) (uint32, error) {
	article := Article{
		AuthorID: uid,
		Title:    title,
		Content:  content,
		Time:     time,
	}

	tx := setting.DB().WithContext(ctx)
	result := tx.Create(&article)
	if result.Error != nil {
		return 0, result.Error
	}

	cache.Hset(ctx, "articles", article.ID, article)
	return article.ID, nil
}

func DeleteArticle(ctx context.Context, id, uid uint32) error {
	tx := setting.DB().WithContext(ctx)
	result := tx.Where("author_id = ?", uid).Delete(&Article{}, uint(id))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("not found")
	}

	cache.Hdel(ctx, "articles", id)
	return nil
}

func UpdateArticle(ctx context.Context, id, uid uint32, title, content string) error {
	var article Article
	tx := setting.DB().WithContext(ctx)
	result := tx.First(&article, id)
	if result.Error != nil {
		return result.Error
	}

	if id <= 0 || result.RowsAffected == 0 || article.AuthorID != uid {
		return errors.New("Forbidden")
	}

	article.Title = title
	article.Content = content
	result = tx.Save(&article)
	if result.Error != nil {
		return result.Error
	}

	cache.Hset(ctx, "articles", id, article)
	return nil
}
