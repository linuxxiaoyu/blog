package cache

import (
	"github.com/garyburd/redigo/redis"
	"github.com/linuxxiaoyu/blog/internal/setting"
)

func Zadd(key string, id uint32, name string) error {
	c := setting.RedisConn()
	defer c.Close()

	_, err := c.Do("ZADD", key, id, name)
	return err
}

func Zscore(key string, name string) (uint32, error) {
	c := setting.RedisConn()
	defer c.Close()

	id, err := redis.Uint64(c.Do("ZSCORE", key, name))
	return uint32(id), err
}
