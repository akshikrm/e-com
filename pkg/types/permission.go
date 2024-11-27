package types

import "time"

type CreateNewPermission struct {
	RoleCode     string `json:"rolecode"`
	ResourceCode string `json:"resourcecode"`
	R            bool   `json:"r"`
	W            bool   `json:"w"`
	U            bool   `json:"u"`
	D            bool   `json:"d"`
}

type Permission struct {
	ID           int       `json:"id"`
	RoleCode     int       `json:"role_code"`
	ResourceCode int       `json:"resource_code"`
	R            bool      `json:"r"`
	W            bool      `json:"w"`
	U            bool      `json:"u"`
	D            bool      `json:"d"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}
