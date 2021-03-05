package data

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"gorm.io/gorm"

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

func GetUsers(ctx context.Context, ids []uint32) (map[uint32]string, error) {
	if ids == nil || len(ids) == 0 {
		return nil, errors.New("ids can't empty")
	}
	tx := setting.DB().WithContext(ctx)
	users := []User{}
	m := map[uint32]string{}
	result := tx.Where("id IN ?", ids).FindInBatches(&users, 100, func(tx *gorm.DB, batch int) error {
		for _, user := range users {
			m[user.ID] = user.Name
		}
		users = users[0:0]
		return nil
	})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected <= 0 {
		return nil, errors.New("not found")
	}
	return m, nil
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

func CreateUser(ctx context.Context, user *User) (uint32, error) {
	if user == nil {
		return 0, errors.New("user is nil")
	}
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
