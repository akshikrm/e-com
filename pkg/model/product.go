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

func (p *ProductModel) Create(product *types.CreateNewProduct) (uint, error) {
	query := `INSERT INTO products(name, slug, price, image, description) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	row := p.store.QueryRow(query,
		product.Name,
		product.Slug,
		product.Price,
		product.Image,
		product.Description,
	)

	savedProduct := Product{}
	if err := row.Scan(&savedProduct.ID); err != nil {
		log.Printf("failed to create new product %s due to %s", product.Name, err)
		return 0, utils.ServerError
	}
	return savedProduct.ID, nil
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
	products := []*Product{}
	for rows.Next() {
		product, err := scanProductRows(rows)
		if err != nil {
			return nil, utils.ServerError
		}
		products = append(products, product)
	}
	return products, nil
}

func scanProductRows(rows *sql.Rows) (*Product, error) {
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
		log.Printf("scan into product failed due to %s", err)
	}

	return &product, err
}

func NewProductModel(store *sql.DB) *ProductModel {
	return &ProductModel{
		store: store,
	}
}