package user

import "gorm.io/gorm"

type Repository interface {
	CreateUser(user *User) error
}

type repository struct {
	db *gorm.DB
}
