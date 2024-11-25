package model

import (
	"akshidas/e-com/pkg/types"
	"akshidas/e-com/pkg/utils"
	"database/sql"
	"log"
	"time"
)

type ProductCategory struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Description string     `json:"description"`
	Enabled     bool       `json:"enabled"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type ProductCategoriesModel struct {
	store *sql.DB
}

func (p *ProductCategoriesModel) Create(newCategory *types.NewProductCategoryRequest) (*ProductCategory, error) {
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

func (p *ProductCategoriesModel) GetAll() ([]*ProductCategory, error) {
	query := "SELECT * FROM  product_categories"
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

func (p *ProductCategoriesModel) GetOne(id int) (*ProductCategory, error) {
	query := "SELECT * FROM product_categories WHERE id=$1"
	row := p.store.QueryRow(query, id)
	productCategory, err := scanCategoryRow(row)
	if err != nil {
		log.Printf("failed to get product category with %s due to %s", id, err)
		if err == sql.ErrNoRows {
			return nil, utils.NotFound
		}
		return nil, utils.ServerError
	}
	return productCategory, err
}

func (p *ProductCategoriesModel) Update(id int, updateProductCategory *types.UpdateProductCategoryRequest) (*ProductCategory, error) {
	query := "UPDATE product_categories SET name=$1, slug=$2, description=$3, enabled=$4 WHERE id=$5 RETURNING *"
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

func (p *ProductCategoriesModel) Delete(id int) error {
	query := "DELETE FROM product_categories WHERE id=$1"

	if _, err := p.store.Exec(query, id); err != nil {
		log.Default().Printf("failed to delete product category with id %d due to %s", id, err)
		return utils.ServerError
	}
	return nil
}

func scanCategoryRow(row *sql.Row) (*ProductCategory, error) {
	productCategory := &ProductCategory{}
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

func scanCategoryRows(rows *sql.Rows) ([]*ProductCategory, error) {
	productsCategories := []*ProductCategory{}

	for rows.Next() {
		productCategory := &ProductCategory{}
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

func NewProductCategories(store *sql.DB) *ProductCategoriesModel {
	return &ProductCategoriesModel{store: store}
}
