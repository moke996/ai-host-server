package global

import (
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

var RedisClint *redis.Client

func InitRedis() {
	redisOptions := &redis.Options{
		Addr:         Config.Redis.Address,
		PoolSize:     Config.Redis.MaxPoolSize,
		MinIdleConns: 1,
		IdleTimeout:  30 * time.Second,
	}

	RedisClint = redis.NewClient(redisOptions)
	log.Println("Redis is Collection!!!")
}

func StopRedis() {
	err := RedisClint.Close()
	if err != nil {
		log.Println(err)
	}
}
