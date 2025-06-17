package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       int
	Username string `gorm:"unique"`
	Password string
	Email    string `gorm:"unique"`
	Source   string `gorm:"source"`
}

func CreateUser(userData *User) error {
	return nil
}

func GetUserByName(username string) (*User, error) {
	return nil, nil
}
