package cache

import (
	"context"

	"github.com/gomodule/redigo/redis"
	"github.com/linuxxiaoyu/blog/internal/setting"
)

func Sadd(ctx context.Context, key, member string) error {
	c, err := setting.RedisConnWithContext(ctx)
	if err != nil {
		return err
	}
	defer c.Close()

	_, err = c.Do("SADD", key, member)
	return err
}

func Srem(key, member string) error {
	c := setting.RedisConn()
	defer c.Close()

	_, err := c.Do("SREM", key, member)
	return err
}

func Smembers(key string) ([]int64, error) {
	c := setting.RedisConn()
	defer c.Close()

	return redis.Int64s(c.Do("SMEMBERS", key))
}
