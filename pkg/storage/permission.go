package storage

import (
	"akshidas/e-com/pkg/types"
	"akshidas/e-com/pkg/utils"
	"database/sql"
	"log"
	"time"
)

type PermissionStorage struct {
	store *sql.DB
}

func (r *PermissionStorage) GetAll() ([]*types.Permission, error) {
	query := `select * from permissions`
	rows, err := r.store.Query(query)

	if err == sql.ErrNoRows {
		return nil, utils.NotFound
	}
	if err != nil {
		log.Printf("Failed to query db for permissions due to %s", err)
		return nil, utils.ServerError
	}

	permissions, err := scanPermissionRows(rows)
	return permissions, err
}

func (r *PermissionStorage) GetOne(id int) (*types.Permission, error) {
	query := `select * from permissions where id=$1`
	row := r.store.QueryRow(query, id)

	permission, err := scanPermissionRow(row)

	if err != nil {
		log.Printf("Failed to query db for permission %d due to %s", id, err)
	}

	return permission, nil
}

func (r *PermissionStorage) Create(newPermission *types.CreateNewPermission) error {
	query := "INSERT INTO permissions(role_code, resource_code, r, w, u, d) VALUES ($1, $2, $3, $4, $5, $6)"
	if _, err := r.store.Exec(query,
		newPermission.RoleCode,
		newPermission.ResourceCode,
		newPermission.R,
		newPermission.W,
		newPermission.U,
		newPermission.D,
	); err != nil {
		return utils.ServerError
	}

	return nil
}

func (r *PermissionStorage) Update(id int, updatedPermission *types.CreateNewPermission) (*types.Permission, error) {
	query := `UPDATE roles SET role_code=$1, resource_code=$2, r=$3, w=$4, u=$5, d=$6 returning *`
	row := r.store.QueryRow(query,
		updatedPermission.RoleCode,
		updatedPermission.ResourceCode,
		updatedPermission.R,
		updatedPermission.W,
		updatedPermission.U,
		updatedPermission.D,
	)

	role, err := scanPermissionRow(row)
	if err != nil {
		log.Printf("failed to update role %d due to %s", id, err)
	}

	return role, nil
}

func (r *PermissionStorage) Delete(id int) error {
	query := "UPDATE roles set deleted_at=$1 where id=$2"
	if _, err := r.store.Exec(query, time.Now(), id); err != nil {
		log.Printf("failed to delete permission %d due to %s", id, err)
		return utils.ServerError
	}
	return nil
}

func scanPermissionRows(rows *sql.Rows) ([]*types.Permission, error) {
	permissions := []*types.Permission{}
	for rows.Next() {
		permission := &types.Permission{}
		err := rows.Scan(
			&permission.ID,
			&permission.RoleCode,
			&permission.ResourceCode,
			&permission.R,
			&permission.W,
			&permission.U,
			&permission.D,
			&permission.CreatedAt,
			&permission.UpdatedAt,
			&permission.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}

	return permissions, nil
}

func scanPermissionRow(row *sql.Row) (*types.Permission, error) {
	permission := types.Permission{}
	err := row.Scan(
		&permission.ID,
		&permission.RoleCode,
		&permission.ResourceCode,
		&permission.R,
		&permission.W,
		&permission.U,
		&permission.D,
		&permission.CreatedAt,
		&permission.UpdatedAt,
		&permission.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, utils.NotFound
	}

	if err != nil {
		return nil, utils.ServerError
	}

	return &permission, nil

}

func NewPermissionStorage(store *sql.DB) *PermissionStorage {
	return &PermissionStorage{
		store: store,
	}
}
