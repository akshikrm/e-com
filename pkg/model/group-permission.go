package model

import (
	"akshidas/e-com/pkg/types"
	"akshidas/e-com/pkg/utils"
	"database/sql"
	"log"
	"time"
)

type GroupPermission struct {
	ID          int       `json:"id"`
	RoleID      int       `json:"role_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GroupPermissionModel struct {
	store *sql.DB
}

func (r *GroupPermissionModel) GetAll() ([]*GroupPermission, error) {
	query := `select * from group_permissions`
	rows, err := r.store.Query(query)

	if err == sql.ErrNoRows {
		return nil, utils.NotFound
	}
	if err != nil {
		log.Printf("Failed to query db for group_permissions due to %s", err)
		return nil, utils.ServerError
	}

	groupPermissions, err := scanGroupPermissionRows(rows)
	return groupPermissions, err
}

func (r *GroupPermissionModel) GetOne(id int) (*GroupPermission, error) {
	query := `select * from group_permissions where id=$1`
	row := r.store.QueryRow(query, id)

	groupPermissions, err := scanGroupPermissionRow(row)

	if err != nil {
		log.Printf("Failed to query db for group_permissions %d due to %s", id, err)
	}

	return groupPermissions, nil
}

func (r *GroupPermissionModel) Create(newGroupPermission *types.CreateNewGroupPermission) error {
	query := "INSERT INTO group_permissions(group_id, permission_id) VALUES ($1, $2)"
	if _, err := r.store.Exec(query,
		newGroupPermission.GroupID,
		newGroupPermission.PermissionID,
	); err != nil {
		return utils.ServerError
	}

	return nil
}

func (r *GroupPermissionModel) Update(id int, updatedGroupPermission *types.CreateNewGroupPermission) (*GroupPermission, error) {
	query := `UPDATE roles SET name=$1, description=$2 role_id=$3 returning *`
	row := r.store.QueryRow(query,
		updatedGroupPermission.GroupID,
		updatedGroupPermission.PermissionID,
	)

	role, err := scanGroupPermissionRow(row)
	if err != nil {
		log.Printf("failed to update role %d due to %s", id, err)
	}

	return role, nil
}

func (r *GroupPermissionModel) Delete(id int) error {
	query := "delete from group_permissions where id=$1"
	_, err := r.store.Exec(query, id)

	if err == sql.ErrNoRows {
		return utils.NotFound
	}
	if err != nil {
		log.Printf("failed to delete group_permissions %d due to %s", id, err)
		return utils.ServerError
	}

	return nil
}

func scanGroupPermissionRows(rows *sql.Rows) ([]*GroupPermission, error) {
	groupPermissions := []*GroupPermission{}
	for rows.Next() {
		groupPermission := &GroupPermission{}
		err := rows.Scan(
			&groupPermission.ID,
			&groupPermission.RoleID,
			&groupPermission.Name,
			&groupPermission.Description,
			&groupPermission.CreatedAt,
			&groupPermission.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		groupPermissions = append(groupPermissions, groupPermission)
	}

	return groupPermissions, nil
}

func scanGroupPermissionRow(row *sql.Row) (*GroupPermission, error) {
	groupPermission := GroupPermission{}
	err := row.Scan(
		&groupPermission.ID,
		&groupPermission.RoleID,
		&groupPermission.Name,
		&groupPermission.Description,
		&groupPermission.CreatedAt,
		&groupPermission.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, utils.NotFound
	}

	if err != nil {
		return nil, utils.ServerError
	}

	return &groupPermission, nil

}

func NewGroupPermissionModel(store *sql.DB) *GroupPermissionModel {
	return &GroupPermissionModel{
		store: store,
	}
}
