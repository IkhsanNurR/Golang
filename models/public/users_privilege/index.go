package users_privilege_model

import (
	"time"
)

type UsersPrivilegeTable struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	IdUsers     int       `json:"id_users" binding:"required"`
	IdPrivilege int       `json:"id_privilege" binding:"required"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime;column:createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime;column:updatedAt"`
	Status      bool      `json:"status" binding:"required"`
}

type UpdateUsersPrivilegeTable struct {
	IdPrivilege int  `json:"id_privilege" binding:"required"`
	Status      bool `json:"status" binding:"required"`
}
