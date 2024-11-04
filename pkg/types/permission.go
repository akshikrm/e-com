package types

type CreateNewPermission struct {
	RoleCode     string `json:"rolecode"`
	ResourceCode string `json:"resourcecode"`
	R            bool   `json:"r"`
	W            bool   `json:"w"`
	U            bool   `json:"u"`
	D            bool   `json:"d"`
}
