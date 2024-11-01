package types

type CreateNewGroup struct {
	Name        string `json:"name"`
	RoleID      int    `json:"role_id"`
	Description string `json:"description"`
}
