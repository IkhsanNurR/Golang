package response

import "time"

type GetAllPrivilege struct {
	ID            int       `json:"id"`
	PrivilegeName string    `json:"privilege_name"`
	Status        bool      `json:"status"`
	CreatedAt     time.Time `json:"createdAt" gorm:"column:createdAt"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"column:updatedAt"`
}
