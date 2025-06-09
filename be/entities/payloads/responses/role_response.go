package responses

import "time"

type (
	RoleListResponse struct {
		ID        uint64     `json:"id"`
		Role      string     `json:"role"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt time.Time  `json:"updated_at"`
		DeletedAt *time.Time `json:"deleted_at"`
	}
)
