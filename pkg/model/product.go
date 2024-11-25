package model

import (
	"akshidas/e-com/pkg/types"
	"akshidas/e-com/pkg/utils"
	"database/sql"
	"log"
	"time"
)

type Product struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Price       uint       `json:"price"`
	Image       string     `json:"image"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
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

func (p *ProductModel) Update(pid int, product *types.CreateNewProduct) (*Product, error) {
	query := `UPDATE products SET name=$1, slug=$2, price=$3, image=$4, description=$5 WHERE id=$6 RETURNING *`
	row := p.store.QueryRow(query,
		product.Name,
		product.Slug,
		product.Price,
		product.Image,
		product.Description,
		pid,
	)

	savedProduct, err := scanProductRow(row)
	if err == sql.ErrNoRows {
		return nil, utils.NotFound
	}

	if err != nil {
		log.Printf("failed to update product %s due to %s", product.Name, err)
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
		log.Printf("failed to get all products due to %s", err)
		return nil, utils.ServerError
	}

	return products, nil
}

func (m *ProductModel) GetOne(id int) (*Product, error) {
	query := `select * from products where id=$1`
	row := m.store.QueryRow(query, id)

	product, err := scanProductRow(row)
	if err == sql.ErrNoRows {

		log.Printf("product with id %d not found due to %s", id, err)
		return nil, utils.NotFound
	}
	if err != nil {
		return nil, utils.ServerError
	}

	return product, nil
}

func (m *ProductModel) Delete(id int) error {
	query := "UPDATE products set deleted_at=$1 where id=$2"
	if _, err := m.store.Exec(query, time.Now(), id); err != nil {
		log.Printf("failed to products %d due to %s", id, err)
		return utils.ServerError
	}
	return nil
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
			&product.DeletedAt,
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
		&product.DeletedAt,
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
