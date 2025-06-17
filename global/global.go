package global

import (
	"github.com/sirupsen/logrus"
)

const (
	MSG_REGISTER   = 1
	MSG_LOGIN      = 2
	MSG_LOGOUT     = 3
	MSG_LOGIN_RES  = 100
	MSG_COMMON_RES = 200
)

var Log = logrus.New()

func init() {
	// RedisPool = &redis.Pool{
	// 	MaxIdle:     5,
	// 	MaxActive:   10,
	// 	IdleTimeout: 300 * time.Second,
	// 	Dial: func() (redis.Conn, error) {
	// 		return redis.Dial("tcp", "localhost:6379")
	// 	},
	// }
	// RedisPool = initialize.InitRedis()
}
