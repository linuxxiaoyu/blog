package data

import (
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

func GetUser(id uint32) (User, error) {
	user, err := getUserFromCache(id)
	if err == nil {
		return user, err
	}

	db := setting.DB()
	result := db.First(&user, id)
	if result.Error != nil {
		return user, result.Error
	}

	cache.Hset("users", id, user)
	cache.Zadd("username", id, user.Name)
	return user, nil
}

func GetUserByName(name string) (User, error) {
	user, err := getUserByNameFromCache(name)
	if err == nil {
		return user, err
	}

	db := setting.DB()
	result := db.Where("name = ?", name).First(&user)
	if result.Error != nil {
		return user, result.Error
	}

	if result.RowsAffected > 0 {
		cache.Zadd("username", user.ID, name)
		cache.Hset("users", user.ID, user)
	}

	return user, nil
}

func getUserFromCache(id uint32) (User, error) {
	var user User
	str, err := cache.Hget("users", id)
	if err == nil {
		err = json.Unmarshal([]byte(str), &user)
	}
	return user, err
}

func getUserByNameFromCache(name string) (User, error) {
	var user User
	id, err := cache.Zscore("username", name)
	if err != nil {
		return user, err
	}
	return getUserFromCache(id)
}

func CreateUser(user User) (uint32, error) {
	if user.Name == "" || user.Password == "" {
		return 0, errors.New("name and password can not empty")
	}

	db := setting.DB()
	result := db.Create(&user)
	if result.Error != nil {
		return 0, errors.WithStack(result.Error)
	}

	cache.Hset("users", user.ID, user)
	cache.Zadd("username", user.ID, user.Name)
	return user.ID, nil
}
