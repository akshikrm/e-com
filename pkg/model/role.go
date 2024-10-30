package model

import (
	"akshidas/e-com/pkg/utils"
	"database/sql"
	"log"
	"time"
)

type Role struct {
	ID          int       `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateRoleRequest struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RoleModel struct {
	store *sql.DB
}

func (r *RoleModel) GetAll() ([]*Role, error) {
	query := `select * from roles`
	rows, err := r.store.Query(query)

	if err == sql.ErrNoRows {
		return nil, utils.NotFound
	}
	if err != nil {
		log.Printf("Failed to query db for roles due to %s", err)
		return nil, utils.ServerError
	}

	roles, err := scanRows(rows)
	return roles, err
}

func (r *RoleModel) GetOne(id int) (*Role, error) {
	query := `select * from roles where id=$1`
	row := r.store.QueryRow(query, id)

	role, err := scanRow(row)

	if err != nil {
		log.Printf("Failed to query db for role %d due to %s", id, err)
	}

	return role, nil
}

func (r *RoleModel) Create(newRole *CreateRoleRequest) error {
	query := `INSERT INTO roles(name, code, description)
	VALUES($1,$2, $3)
	`

	if _, err := r.store.Exec(query,
		newRole.Name,
		newRole.Code,
		newRole.Description,
	); err != nil {
		return utils.ServerError
	}

	return nil
}

func (r *RoleModel) Update(id int, newRole *CreateRoleRequest) (*Role, error) {
	query := `UPDATE roles SET name=$1, code=$2, description=$3 returning *`
	row := r.store.QueryRow(query,
		newRole.Name,
		newRole.Code,
		newRole.Description,
	)

	role, err := scanRow(row)
	if err != nil {
		log.Printf("failed to update role %d due to %s", id, err)
	}

	return role, nil
}

func (r *RoleModel) Delete(id int) error {
	query := "delete from roles where id=$1"
	_, err := r.store.Exec(query, id)

	if err == sql.ErrNoRows {
		return utils.NotFound
	}
	if err != nil {
		log.Printf("failed to delete role %d due to %s", id, err)
		return utils.ServerError
	}

	return nil
}

func scanRows(rows *sql.Rows) ([]*Role, error) {
	roles := []*Role{}
	for rows.Next() {
		role := &Role{}
		err := rows.Scan(
			&role.ID,
			&role.Code,
			&role.Name,
			&role.Description,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		roles = append(roles, role)
	}

	return roles, nil
}

func scanRow(row *sql.Row) (*Role, error) {
	role := Role{}
	err := row.Scan(
		&role.ID,
		&role.Code,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, utils.NotFound
	}

	if err != nil {
		return nil, utils.ServerError
	}

	return &role, nil

}

func NewRoleModel(store *sql.DB) *RoleModel {
	return &RoleModel{
		store: store,
	}
}
