package cache

import (
	"encoding/json"

	"github.com/gomodule/redigo/redis"
	"github.com/linuxxiaoyu/blog/internal/setting"
)

func Set(key string, data interface{}) error {
	c := setting.RedisConn()
	defer c.Close()

	str := ""
	s, ok := data.(string)
	if ok {
		str = s
	} else {
		bs, err := json.Marshal(data)
		if err != nil {
			return err
		}
		str = string(bs)
	}

	_, err := c.Do("SET", key, str)
	return err
}

func Get(key string) (string, error) {
	c := setting.RedisConn()
	defer c.Close()

	return redis.String(c.Do("GET", key))
}

func Del(key string) error {
	c := setting.RedisConn()
	defer c.Close()

	_, err := c.Do("DEL", key)
	return err
}
