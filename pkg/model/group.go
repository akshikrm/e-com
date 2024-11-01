package model

import (
	"akshidas/e-com/pkg/types"
	"akshidas/e-com/pkg/utils"
	"database/sql"
	"log"
	"time"
)

type Group struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GroupModel struct {
	store *sql.DB
}

func (r *GroupModel) GetAll() ([]*Group, error) {
	query := `select * from groups`
	rows, err := r.store.Query(query)

	if err == sql.ErrNoRows {
		return nil, utils.NotFound
	}
	if err != nil {
		log.Printf("Failed to query db for groups due to %s", err)
		return nil, utils.ServerError
	}

	groups, err := scanGroupRows(rows)
	return groups, err
}

func (r *GroupModel) GetOne(id int) (*Group, error) {
	query := `select * from groups where id=$1`
	row := r.store.QueryRow(query, id)

	groups, err := scanGroupRow(row)

	if err != nil {
		log.Printf("Failed to query db for groups %d due to %s", id, err)
	}

	return groups, nil
}

func (r *GroupModel) Create(newGroup *types.CreateNewGroup) error {
	query := "INSERT INTO groups(name, description) VALUES ($1, $2)"
	if _, err := r.store.Exec(query,
		newGroup.Name,
		newGroup.Description,
	); err != nil {
		return utils.ServerError
	}

	return nil
}

func (r *GroupModel) Update(id int, updatedGroup *types.CreateNewGroup) (*Group, error) {
	query := `UPDATE roles SET name=$1, description=$2 returning *`
	row := r.store.QueryRow(query,
		updatedGroup.Name,
		updatedGroup.Description,
	)

	role, err := scanGroupRow(row)
	if err != nil {
		log.Printf("failed to update role %d due to %s", id, err)
	}

	return role, nil
}

func (r *GroupModel) Delete(id int) error {
	query := "delete from groups where id=$1"
	_, err := r.store.Exec(query, id)

	if err == sql.ErrNoRows {
		return utils.NotFound
	}
	if err != nil {
		log.Printf("failed to delete groups %d due to %s", id, err)
		return utils.ServerError
	}

	return nil
}

func scanGroupRows(rows *sql.Rows) ([]*Group, error) {
	groups := []*Group{}
	for rows.Next() {
		group := &Group{}
		err := rows.Scan(
			&group.ID,
			&group.CreatedAt,
			&group.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		groups = append(groups, group)
	}

	return groups, nil
}

func scanGroupRow(row *sql.Row) (*Group, error) {
	groups := Group{}
	err := row.Scan(
		&groups.ID,
		&groups.CreatedAt,
		&groups.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, utils.NotFound
	}

	if err != nil {
		return nil, utils.ServerError
	}

	return &groups, nil

}

func NewGroupModel(store *sql.DB) *GroupModel {
	return &GroupModel{
		store: store,
	}
}
