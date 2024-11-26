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
	CategoryID  uint       `json:"category_id"`
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
	query := `INSERT INTO products(name, slug, price, image, description, category_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`
	row := p.store.QueryRow(query,
		product.Name,
		product.Slug,
		product.Price,
		product.Image,
		product.Description,
		product.CategoryID,
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
	query := "UPDATE products SET name=$1, slug=$2, price=$3, image=$4, description=$5, category_id=$6 WHERE id=$7 AND deleted_at IS NULL RETURNING *"
	row := p.store.QueryRow(query,
		product.Name,
		product.Slug,
		product.Price,
		product.Image,
		product.Description,
		product.CategoryID,
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

func (p *ProductModel) GetAll() ([]*types.ProductsList, error) {
	query := "SELECT p.id, p.name, p.slug, p.price, p.image, p.description, c.id as c_id, c.name as c_name,c.slug as c_slug,c.description as c_description FROM products p INNER JOIN product_categories c ON p.category_id=c.id AND c.enabled='t' where p.deleted_at IS NULL;"
	rows, err := p.store.Query(query)

	if err == sql.ErrNoRows {
		return nil, utils.NotFound
	}

	if err != nil {
		log.Printf("failed to get all products due to %s", err)
		return nil, utils.ServerError
	}

	products := []*types.ProductsList{}
	for rows.Next() {
		product := types.ProductsList{}
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Slug,
			&product.Price,
			&product.Image,
			&product.Description,
			&product.Category.ID,
			&product.Category.Name,
			&product.Category.Slug,
			&product.Category.Description,
		)
		if err != nil {
			log.Printf("failed to scan products due to %s", err)
			return nil, utils.ServerError
		}
		products = append(products, &product)
	}
	return products, nil
}

func (m *ProductModel) GetOne(id int) (*Product, error) {
	query := "SELECT * FROM products WHERE id=$1 AND deleted_at IS NULL"
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
	query := "UPDATE products SET deleted_at=$1 WHERE id=$2"
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
			&product.CategoryID,
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
		&product.CategoryID,
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
