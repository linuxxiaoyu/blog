package cache

import (
	"encoding/json"

	"github.com/gomodule/redigo/redis"
	"github.com/linuxxiaoyu/blog/pkg/setting"
)

func Hset(key string, id uint, value interface{}) error {
	c := setting.RedisConn()
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

	_, err := c.Do("HSET", key, id, str)
	return err
}

func Hget(key string, id uint) (string, error) {
	c := setting.RedisConn()
	defer c.Close()

	return redis.String(c.Do("HGET", key, id))
}

func Hdel(key string, id uint) error {
	c := setting.RedisConn()
	defer c.Close()

	_, err := c.Do("HDEL", key, id)
	return err
}
