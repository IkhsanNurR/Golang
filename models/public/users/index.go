package users_model

import (
	"time"
)

type UserTable struct {
	ID             int       `json:"id" gorm:"primaryKey"`
	FullName       string    `json:"full_name" binding:"required"`
	Email          string    `json:"email" binding:"required"`
	Password       string    `json:"password" binding:"required"`
	Status         bool      `json:"status" binding:"required"`
	ProfilePicture string    `json:"profile_picture" gorm:"default:null ;column:profile_picture"`
	CreatedAt      time.Time `json:"createdAt" gorm:"autoCreateTime;column:createdAt"`
	UpdatedAt      time.Time `json:"updatedAt" gorm:"autoUpdateTime;column:updatedAt"`
}

type UpdateUserTable struct {
	FullName       string    `json:"full_name" binding:"required"`
	Email          string    `json:"email" binding:"required"`
	Status         bool      `json:"status" binding:"required"`
	ProfilePicture string    `json:"profile_picture" gorm:"default:null ;column:profile_picture"`
	CreatedAt      time.Time `json:"createdAt" gorm:"autoCreateTime;column:createdAt"`
	UpdatedAt      time.Time `json:"updatedAt" gorm:"autoUpdateTime;column:updatedAt"`
}
