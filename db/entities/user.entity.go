package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserPubId string `gorm:"unique"`
	Username  *string
	IsPublic  bool `gorm:"default=false"`
}
