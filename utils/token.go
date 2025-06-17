package utils

import (
	"crypto/rand"
	"encoding/base64"

	// "game_server/global"
	"game_server/initialize"

	"github.com/garyburd/redigo/redis"
)

// 生成Token
func GenerateToken(userID int) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(b)

	// 存储到Redis (有效期7天)
	conn := initialize.RedisPool.Get()
	defer conn.Close()

	_, err := conn.Do("SETEX", "token:"+token, 7*24*3600, userID)
	return token, err
}

// 验证Token
func VerifyToken(token string) (int, bool) {
	conn := initialize.RedisPool.Get()
	defer conn.Close()

	userID, err := redis.Int(conn.Do("GET", "token:"+token))
	if err != nil {
		return 0, false
	}

	// 刷新Token有效期
	conn.Do("EXPIRE", "token:"+token, 7*24*3600)
	return userID, true
}
