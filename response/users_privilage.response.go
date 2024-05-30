package response

import "time"

type GetAllUsersPrivilege struct {
	ID                   *int       `json:"id"`
	IdUsers              *int       `json:"id_users"`
	IdPrivilege          *int       `json:"id_privilege"`
	StatusUsersPrivilege *bool      `json:"status_users_privilege "`
	CreatedAt            *time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt            *time.Time `json:"updatedAt" gorm:"column:updated_at"`
	FullName             *string    `json:"full_name"`
	Email                *string    `json:"email"`
	StatusUsers          *bool      `json:"status_users"`
	PrivilegeName        *string    `json:"privilege_name"`
	StatusPrivilege      *bool      `json:"status_privilege"`
}
