package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name"  gorm:"type:varchar(100);not null"`
	Email    string `json:"email"  gorm:"type:varchar(255);uniqueIndex;not nll"`
	Password string `json:"password"  gorm:"type:varchar(100);not null"`
}
