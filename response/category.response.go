package response

import "time"

type GetAllCategory struct {
	ID             int       `json:"id" gorm:"primaryKey"`
	CategoryName   string    `json:"category_name"`
	CategoryDetail string    `json:"category_detail"`
	CreatedAt      time.Time `json:"createdAt" gorm:"autoCreateTime;column:createdAt"`
	UpdatedAt      time.Time `json:"updatedAt" gorm:"autoUpdateTime;column:updatedAt"`
	Status         bool      `json:"status"`
}
