package privilege

import "time"

type PrivilegeTable struct {
	ID            int       `json:"id" gorm:"primaryKey"`
	PrivilegeName string    `json:"privilege_name" binding:"required"`
	CreatedAt     time.Time `json:"createdAt" gorm:"autoCreateTime;column:createdAt"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"autoUpdateTime;column:updatedAt"`
	Status        bool      `json:"status" binding:"required"`
}
