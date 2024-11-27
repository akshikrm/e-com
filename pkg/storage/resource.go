package storage

import (
	"akshidas/e-com/pkg/types"
	"akshidas/e-com/pkg/utils"
	"database/sql"
	"log"
	"time"
)

type ResourceStorage struct {
	store *sql.DB
}

func (r *ResourceStorage) GetAll() ([]*types.Resource, error) {
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

func (r *ResourceStorage) GetOne(id int) (*types.Resource, error) {
	query := `select * from resources where id=$1`
	row := r.store.QueryRow(query, id)

	resource, err := scanResourceRow(row)

	if err != nil {
		log.Printf("Failed to query db for resource %d due to %s", id, err)
	}

	return resource, nil
}

func (r *ResourceStorage) Create(newResource *types.CreateResourceRequest) error {
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

func (r *ResourceStorage) Update(id int, newResource *types.CreateResourceRequest) (*types.Resource, error) {
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

func (r *ResourceStorage) Delete(id int) error {
	query := "UPDATE resources set deleted_at=$1 where id=$2"
	if _, err := r.store.Exec(query, time.Now(), id); err != nil {
		log.Printf("failed to delete resource %d due to %s", id, err)
		return utils.ServerError
	}
	return nil
}

func scanResourceRows(rows *sql.Rows) ([]*types.Resource, error) {
	resources := []*types.Resource{}
	for rows.Next() {
		resource := &types.Resource{}
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

func scanResourceRow(row *sql.Row) (*types.Resource, error) {
	resource := types.Resource{}
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

func NewResourceStorage(store *sql.DB) *ResourceStorage {
	return &ResourceStorage{
		store: store,
	}
}
