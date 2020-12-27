package data

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/linuxxiaoyu/blog/internal/cache"
	"github.com/linuxxiaoyu/blog/internal/setting"
)

// User is a table in db
type User struct {
	ID       uint32 `gorm:"primary_key"`
	Name     string
	Password string
}

func GetUser(ctx context.Context, id uint32) (User, error) {
	user, err := getUserFromCache(ctx, id)
	if err == nil {
		return user, err
	}

	tx := setting.DB().WithContext(ctx)
	result := tx.First(&user, id)
	if result.Error != nil {
		return user, result.Error
	}

	cache.Hset(ctx, "users", id, user)
	cache.Zadd(ctx, "username", id, user.Name)
	return user, nil
}

func GetUserByName(ctx context.Context, name string) (User, error) {
	user, err := getUserByNameFromCache(ctx, name)
	if err == nil {
		return user, err
	}

	tx := setting.DB().WithContext(ctx)
	result := tx.Where("name = ?", name).First(&user)
	if result.Error != nil {
		return user, result.Error
	}

	if result.RowsAffected > 0 {
		cache.Zadd(ctx, "username", user.ID, name)
		cache.Hset(ctx, "users", user.ID, user)
	}

	return user, nil
}

func getUserFromCache(ctx context.Context, id uint32) (User, error) {
	var user User
	str, err := cache.Hget(ctx, "users", id)
	if err == nil {
		err = json.Unmarshal([]byte(str), &user)
	}
	return user, err
}

func getUserByNameFromCache(ctx context.Context, name string) (User, error) {
	var user User
	id, err := cache.Zscore(ctx, "username", name)
	if err != nil {
		return user, err
	}
	return getUserFromCache(ctx, id)
}

func CreateUser(ctx context.Context, user User) (uint32, error) {
	if user.Name == "" || user.Password == "" {
		return 0, errors.New("name and password can not empty")
	}

	db := setting.DB().WithContext(ctx)
	result := db.Create(&user)
	if result.Error != nil {
		return 0, errors.WithStack(result.Error)
	}

	cache.Hset(ctx, "users", user.ID, user)
	cache.Zadd(ctx, "username", user.ID, user.Name)
	return user.ID, nil
}
