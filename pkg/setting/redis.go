package setting

import (
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
)

var (
	pool *redis.Pool
)

func initCache() {
	redisSection, err := cfg.GetSection("redis")
	if err != nil {
		log.Fatalf("Fail to get section 'redis': %v", err)
	}

	host := redisSection.Key("Host").MustString("127.0.0.1")
	port := redisSection.Key("Port").MustUint(6379)
	// user := redisSection.Key("User").MustString("")
	// password := redisSection.Key("Password").MustString("")
	maxIdle := redisSection.Key("MaxIdle").MustUint(30)

	pool = redis.NewPool(func() (redis.Conn, error) {
		return redis.Dial(
			"tcp",
			fmt.Sprintf("%s:%d", host, port),
			// redis.DialPassword(password),
			// redis.DialUser
		)
	}, int(maxIdle))
}

// RedisConn returns a Conn from redis pool
func RedisConn() redis.Conn {
	return pool.Get()
}
