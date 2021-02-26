package cache

import (
	"context"
	"encoding/json"

	"github.com/gomodule/redigo/redis"
	"github.com/linuxxiaoyu/blog/internal/setting"
)

func Hset(ctx context.Context, key string, id uint32, value interface{}) error {
	c, err := setting.RedisConnWithContext(ctx)
	if err != nil {
		return err
	}
	defer c.Close()

	str := ""
	v, ok := value.(string)
	if ok {
		str = v
	} else {
		bs, err := json.Marshal(value)
		if err != nil {
			return err
		}
		str = string(bs)
	}

	_, err = c.Do("HSET", key, id, str)
	return err
}

func Hget(ctx context.Context, key string, id uint32) (string, error) {
	c, err := setting.RedisConnWithContext(ctx)
	if err != nil {
		return "", err
	}
	defer c.Close()

	return redis.String(c.Do("HGET", key, id))
}

func Hdel(ctx context.Context, key string, id uint32) error {
	c, err := setting.RedisConnWithContext(ctx)
	defer c.Close()

	_, err = c.Do("HDEL", key, id)
	return err
}
