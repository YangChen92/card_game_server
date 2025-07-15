package initialize

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

var RedisPool *redis.Pool

func InitConfig() {

}

func InitMySQL() {
}

func InitRedis() {
	// 初始化redis连接池
	RedisPool = &redis.Pool{
		MaxIdle:     5,
		MaxActive:   10,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}
	_, err := RedisPool.Dial()
	if err != nil {
		panic(err)
	}

}
