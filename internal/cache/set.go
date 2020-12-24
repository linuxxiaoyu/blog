package cache

import (
	"github.com/gomodule/redigo/redis"
	"github.com/linuxxiaoyu/blog/internal/setting"
)

func Sadd(key, member string) error {
	c := setting.RedisConn()
	defer c.Close()

	_, err := c.Do("SADD", key, member)
	return err
}

func Srem(key, member string) error {
	c := setting.RedisConn()
	defer c.Close()

	_, err := c.Do("SREM", key, member)
	return err
}

func Smembers(key string) ([]string, error) {
	c := setting.RedisConn()
	defer c.Close()

	return redis.Strings(c.Do("SMEMBERS", key))
}
