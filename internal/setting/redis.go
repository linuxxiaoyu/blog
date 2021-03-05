package setting

import (
	"context"
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
)

var (
	pool *redis.Pool
)

func InitCache() {
	initCfg()
	redisSection, err := cfg.GetSection("redis")
	if err != nil {
		log.Fatalf("Fail to get section 'redis': %v", err)
	}

	host := redisSection.Key("Host").MustString("127.0.0.1")
	port := redisSection.Key("Port").MustUint(6379)
	password := redisSection.Key("Password").MustString("")
	maxIdle := redisSection.Key("MaxIdle").MustUint(30)

	pool = redis.NewPool(func() (redis.Conn, error) {
		return redis.Dial(
			"tcp",
			fmt.Sprintf("%s:%d", host, port),
			redis.DialPassword(password),
		)
	}, int(maxIdle))
}

// RedisConn returns a Conn from redis pool
func RedisConn() redis.Conn {
	if pool == nil {
		InitCache()
	}
	return pool.Get()
}

func RedisConnWithContext(ctx context.Context) (redis.Conn, error) {
	if pool == nil {
		InitCache()
	}
	return pool.GetContext(ctx)
}
