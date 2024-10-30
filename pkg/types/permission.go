package types

type CreateNewPermission struct {
	RoleCode     int  `json:"rolecode"`
	ResourceCode int  `json:"resourcecode"`
	R            bool `json:"r"`
	W            bool `json:"w"`
	U            bool `json:"u"`
	D            bool `json:"d"`
}
