package cache

import (
	"context"

	"github.com/garyburd/redigo/redis"
	"github.com/linuxxiaoyu/blog/internal/setting"
)

func Zadd(ctx context.Context, key string, id uint32, name string) error {
	c, err := setting.RedisConnWithContext(ctx)
	if err != nil {
		return err
	}
	defer c.Close()

	_, err = c.Do("ZADD", key, id, name)
	return err
}

func Zscore(ctx context.Context, key string, name string) (uint32, error) {
	c, err := setting.RedisConnWithContext(ctx)
	if err != nil {
		return 0, err
	}
	defer c.Close()

	id, err := redis.Uint64(c.Do("ZSCORE", key, name))
	return uint32(id), err
}
