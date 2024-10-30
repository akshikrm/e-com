package model

// There will be a table which connects users to groups with  group_id and user_id columns

type Group struct {
	ID          int
	Name        string
	Description string
	Roles       int // this will be accomplished by a table which connects groups and roles
}

type Resource struct {
	ID          int
	Code        string
	Name        string
	Description string
}

type RXWD struct {
	R bool // Read
	W bool // Write
	X bool // Execute
	D bool // Delete
}

type RoleResourcePermission struct {
	ID           int
	RoleCode     int
	ResourceCode int
	Permission   RXWD
}
