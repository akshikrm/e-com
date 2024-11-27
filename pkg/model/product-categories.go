package model

import (
	"akshidas/e-com/pkg/types"
	"akshidas/e-com/pkg/utils"
	"database/sql"
	"log"
	"time"
)

type ProductCategoriesStorage struct {
	store *sql.DB
}

func (p *ProductCategoriesStorage) Create(newCategory *types.NewProductCategoryRequest) (*types.ProductCategory, error) {
	query := "INSERT INTO product_categories(name, slug, description, enabled) VALUES($1, $2, $3, $4) RETURNING *"
	row := p.store.QueryRow(query, newCategory.Name, newCategory.Slug, newCategory.Description, newCategory.Enabled)

	savedCategory, err := scanCategoryRow(row)
	if err == sql.ErrNoRows {
		return nil, utils.NotFound
	}

	if err != nil {
		log.Printf("Failed to create category %s due to %s", newCategory.Name, err)
		return nil, utils.ServerError
	}

	return savedCategory, nil

}

func (p *ProductCategoriesStorage) GetAll() ([]*types.ProductCategory, error) {
	query := "SELECT * FROM  product_categories WHERE deleted_at IS NULL"
	rows, err := p.store.Query(query)
	if err == sql.ErrNoRows {
		return nil, utils.NotFound
	}
	if err != nil {
		log.Printf("failed to get product_categories due to %s", err)
		return nil, utils.ServerError
	}
	productsCategories, err := scanCategoryRows(rows)
	if err != nil {
		log.Printf("failed to get all products due to %s", err)
		return nil, utils.ServerError
	}
	return productsCategories, nil
}

func (p *ProductCategoriesStorage) GetOne(id int) (*types.ProductCategory, error) {
	query := "SELECT * FROM product_categories WHERE id=$1 AND deleted_at IS NULL"
	row := p.store.QueryRow(query, id)
	productCategory, err := scanCategoryRow(row)
	if err != nil {
		log.Printf("failed to get product category with %d due to %s", id, err)
		if err == sql.ErrNoRows {
			return nil, utils.NotFound
		}
		return nil, utils.ServerError
	}
	return productCategory, err
}

func (p *ProductCategoriesStorage) Update(id int, updateProductCategory *types.UpdateProductCategoryRequest) (*types.ProductCategory, error) {
	query := "UPDATE product_categories SET name=$1, slug=$2, description=$3, enabled=$4 WHERE id=$5 AND deleted_at IS NULL RETURNING *"
	row := p.store.QueryRow(
		query,
		updateProductCategory.Name,
		updateProductCategory.Slug,
		updateProductCategory.Description,
		updateProductCategory.Enabled,
		id,
	)
	productCategory, err := scanCategoryRow(row)
	if err != nil {
		log.Printf("failed to update product category with id %d due to %s", id, err)
		if err == sql.ErrNoRows {
			return nil, utils.NotFound
		}
		return nil, utils.ServerError
	}
	return productCategory, err
}

func (p *ProductCategoriesStorage) Delete(id int) error {
	query := "UPDATE product_categories set deleted_at=$1 where id=$2 AND deleted_at IS NULL"
	if _, err := p.store.Exec(query, time.Now(), id); err != nil {
		log.Printf("failed to delete product category with id %d due to %s", id, err)
		return utils.ServerError
	}
	return nil
}

func scanCategoryRow(row *sql.Row) (*types.ProductCategory, error) {
	productCategory := &types.ProductCategory{}
	err := row.Scan(
		&productCategory.ID,
		&productCategory.Name,
		&productCategory.Slug,
		&productCategory.Enabled,
		&productCategory.Description,
		&productCategory.CreatedAt,
		&productCategory.UpdatedAt,
		&productCategory.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return productCategory, nil
}

func scanCategoryRows(rows *sql.Rows) ([]*types.ProductCategory, error) {
	productsCategories := []*types.ProductCategory{}

	for rows.Next() {
		productCategory := &types.ProductCategory{}
		err := rows.Scan(
			&productCategory.ID,
			&productCategory.Name,
			&productCategory.Slug,
			&productCategory.Enabled,
			&productCategory.Description,
			&productCategory.CreatedAt,
			&productCategory.UpdatedAt,
			&productCategory.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		productsCategories = append(productsCategories, productCategory)
	}
	return productsCategories, nil
}

func NewProductCategoryStorage(store *sql.DB) *ProductCategoriesStorage {
	return &ProductCategoriesStorage{store: store}
}
