package model

import (
	"akshidas/e-com/pkg/types"
	"akshidas/e-com/pkg/utils"
	"database/sql"
	"log"
	"time"
)

type Product struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Price       uint      `json:"price"`
	Image       string    `json:"image"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductModel struct {
	store *sql.DB
}

func (p *ProductModel) Create(product *types.CreateNewProduct) (*Product, error) {
	query := `INSERT INTO products(name, slug, price, image, description) VALUES ($1, $2, $3, $4, $5) RETURNING *`
	row := p.store.QueryRow(query,
		product.Name,
		product.Slug,
		product.Price,
		product.Image,
		product.Description,
	)

	savedProduct, err := scanProductRow(row)
	if err == sql.ErrNoRows {
		return nil, utils.NotFound
	}
	if err != nil {
		log.Printf("failed to create new product %s due to %s", product.Name, err)
		return nil, utils.ServerError
	}
	return savedProduct, nil
}

func (p *ProductModel) GetAll() ([]*Product, error) {
	query := "SELECT * FROM products;"
	rows, err := p.store.Query(query)

	if err == sql.ErrNoRows {
		return nil, utils.NotFound
	}

	if err != nil {
		log.Printf("failed to get all products due to %s", err)
		return nil, utils.ServerError
	}

	products, err := scanProductRows(rows)
	if err != nil {
		log.Printf("Failed to save product")
		return nil, utils.ServerError
	}

	return products, nil
}

func scanProductRows(rows *sql.Rows) ([]*Product, error) {
	products := []*Product{}
	for rows.Next() {
		product := Product{}
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Slug,
			&product.Price,
			&product.Image,
			&product.Description,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}
	return products, nil
}

func scanProductRow(rows *sql.Row) (*Product, error) {
	product := Product{}
	err := rows.Scan(
		&product.ID,
		&product.Name,
		&product.Slug,
		&product.Price,
		&product.Image,
		&product.Description,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func NewProductModel(store *sql.DB) *ProductModel {
	return &ProductModel{
		store: store,
	}
}
