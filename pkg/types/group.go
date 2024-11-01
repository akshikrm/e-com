package types

type CreateNewGroup struct {
	Name        string `json:"name"`
	RoleID      int    `json:"role_id"`
	Description string `json:"description"`
}

type CreateNewGroupPermission struct {
	GroupID      int    `json:"group_id"`
	PermissionID string `json:"permission_id"`
}
