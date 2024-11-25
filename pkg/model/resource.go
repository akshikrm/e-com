package model

import (
	"akshidas/e-com/pkg/types"
	"akshidas/e-com/pkg/utils"
	"database/sql"
	"log"
	"time"
)

type Resource struct {
	ID          int        `json:"id"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type ResourceModel struct {
	store *sql.DB
}

func (r *ResourceModel) GetAll() ([]*Resource, error) {
	query := `select * from resources`
	rows, err := r.store.Query(query)

	if err == sql.ErrNoRows {
		return nil, utils.NotFound
	}
	if err != nil {
		log.Printf("Failed to query db for resources due to %s", err)
		return nil, utils.ServerError
	}

	resources, err := scanResourceRows(rows)
	return resources, err
}

func (r *ResourceModel) GetOne(id int) (*Resource, error) {
	query := `select * from resources where id=$1`
	row := r.store.QueryRow(query, id)

	resource, err := scanResourceRow(row)

	if err != nil {
		log.Printf("Failed to query db for resource %d due to %s", id, err)
	}

	return resource, nil
}

func (r *ResourceModel) Create(newResource *types.CreateResourceRequest) error {
	query := `INSERT INTO resources(name, code, description)
	VALUES($1,$2, $3)
	`

	if _, err := r.store.Exec(query,
		newResource.Name,
		newResource.Code,
		newResource.Description,
	); err != nil {
		return utils.ServerError
	}

	return nil
}

func (r *ResourceModel) Update(id int, newResource *types.CreateResourceRequest) (*Resource, error) {
	query := `UPDATE resources SET name=$1, code=$2, description=$3 returning *`
	row := r.store.QueryRow(query,
		newResource.Name,
		newResource.Code,
		newResource.Description,
	)

	resource, err := scanResourceRow(row)
	if err != nil {
		log.Printf("failed to update resource %d due to %s", id, err)
	}

	return resource, nil
}

func (r *ResourceModel) Delete(id int) error {
	query := "UPDATE resources set deleted_at=$1 where id=$2"
	if _, err := r.store.Exec(query, time.Now(), id); err != nil {
		log.Printf("failed to delete resource %d due to %s", id, err)
		return utils.ServerError
	}
	return nil
}

func scanResourceRows(rows *sql.Rows) ([]*Resource, error) {
	resources := []*Resource{}
	for rows.Next() {
		resource := &Resource{}
		err := rows.Scan(
			&resource.ID,
			&resource.Code,
			&resource.Name,
			&resource.Description,
			&resource.CreatedAt,
			&resource.UpdatedAt,
			&resource.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		resources = append(resources, resource)
	}

	return resources, nil
}

func scanResourceRow(row *sql.Row) (*Resource, error) {
	resource := Resource{}
	err := row.Scan(
		&resource.ID,
		&resource.Code,
		&resource.Name,
		&resource.Description,
		&resource.CreatedAt,
		&resource.UpdatedAt,
		&resource.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, utils.NotFound
	}

	if err != nil {
		return nil, utils.ServerError
	}

	return &resource, nil

}

func NewResourceModel(store *sql.DB) *ResourceModel {
	return &ResourceModel{
		store: store,
	}
}
