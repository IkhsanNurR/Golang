package category

import "time"

type CategoryTable struct {
	ID             int       `json:"id" gorm:"primaryKey"`
	CategoryName   string    `json:"category_name" binding:"required"`
	CategoryDetail string    `json:"category_detail" binding:"required"`
	CreatedAt      time.Time `json:"createdAt" gorm:"autoCreateTime;column:createdAt"`
	UpdatedAt      time.Time `json:"updatedAt" gorm:"autoUpdateTime;column:updatedAt"`
	Status         bool      `json:"status" binding:"required"`
}
