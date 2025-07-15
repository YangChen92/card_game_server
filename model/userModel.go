package model

import (
	"encoding/json"
	"game_server/initialize"

	"fmt"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID            int    `gorm:"id"`
	Username      string `gorm:"unique"`
	Password      string `gorm:"password"`
	Email         string `gorm:"unique"`
	Source        string `gorm:"source"`
	HeadImg       string `gorm:"head_img"`
	Nickname      string `gorm:"nickname"`
	UserID        int32  `gorm:"user_id"`
	Exp           int32  `gorm:"exp"`
	Money         uint32 `gorm:"money"`
	DeviceID      string `gorm:"device_id"`
	RegTime       string `gorm:"reg_time"`
	LastLoginTime string `gorm:"last_login_time"`
	LastIP        string `gorm:"last_ip"`
}

var db *gorm.DB

// 创建用户
func CreateUser(userData *User) error {
	//先查看redis中是否有数据
	// redisConn := initialize.RedisPool.Get()
	// defer redisConn.Close()
	// key := "user:" + fmt.Sprint(userData.ID)
	// redisConn.Do("GET", key)
	// if redisConn.Err() != nil {
	// 	//如果redis中没有数据，则插入redis
	// 	userData, err := json.Marshal(userData)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	redisConn.Do("SET", key, userData)
	// } else {
	// 	//如果redis中有数据，则直接返回错误
	// 	return fmt.Errorf("用户已存在")
	// }
	userData, err := GetUserData(userData.UserID)
	if err != nil {
		fmt.Println("redis中没有数据 err: ", err)
		return err
	}
	//如果没有数据，则插入数据库
	if err := db.Create(userData).Error; err != nil {
		fmt.Println("插入数据库失败 err: ", err)
		return err
	}
	return nil
}

// 根据用户名获取用户信息
func GetUserByName(username string) (*User, error) {

	var user User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserData(userID int32) (*User, error) {
	redisConn := initialize.RedisPool.Get()
	defer redisConn.Close()
	key := "user:" + fmt.Sprint(userID)
	user, err := redisConn.Do("GET", key)
	if err != nil {
		return nil, err
	}
	if user == nil {
		//如果redis中没有数据，则从数据库中查询
		var user User
		if err := db.Where("user_id = ?", userID).First(&user).Error; err != nil {

			
			return nil, err
		}
		userData, err := json.Marshal(user)
		if err != nil {
			return nil, err
		}
		redisConn.Do("SET", key, userData)
		return &user, nil
	} else {
		//如果redis中有数据，则直接解析
		userData := user.([]byte)
		user := User{}
		if err := json.Unmarshal(userData, &user); err != nil {
			return nil, err
		}
		return &user, nil
	}
}
