package utils

import (
	// "game_server/global"
	"game_server/initialize"

	"github.com/garyburd/redigo/redis"
)

func RegisterLimit(ip string) bool {
	conn := initialize.RedisPool.Get()
	defer conn.Close()

	key := "reg_limit:" + ip
	count, _ := redis.Int(conn.Do("INCR", key))

	if count == 1 {
		conn.Do("EXPIRE", key, 60) // 60秒限制
	}

	return count <= 5 // 每分钟最多5次注册
}
