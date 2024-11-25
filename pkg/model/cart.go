package model

import (
	"akshidas/e-com/pkg/types"
	"akshidas/e-com/pkg/utils"
	"database/sql"
	"log"
	"time"
)

type Cart struct {
	ID        uint       `json:"id"`
	UserID    uint       `json:"user_id"`
	ProductID uint       `json:"product_id"`
	Quantity  uint       `json:"quantity"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type CartModel struct {
	store *sql.DB
}

func (c *CartModel) GetAll(userID uint) ([]*Cart, error) {
	query := "SELECT * FROM carts where user_id=$1"
	rows, err := c.store.Query(query, userID)
	if err == sql.ErrNoRows {
		return nil, utils.NotFound
	}
	if err != nil {
		log.Printf("failed to get all products due to %s", err)
		return nil, utils.ServerError
	}
	return scanCartRows(rows)
}

func (c *CartModel) GetOne(cid uint) (*Cart, error) {
	query := "SELECT * FROM carts WHERE id=$1"
	row := c.store.QueryRow(query, cid)
	cart, err := scanNewCartRow(row)
	if err == sql.ErrNoRows {
		return nil, utils.NotFound
	}
	if err != nil {
		log.Printf("failed to get cart with id %d due to %s", cid, err)
		return nil, utils.ServerError
	}
	return cart, nil
}

func (c *CartModel) Create(newCart *types.CreateCartRequest) (*Cart, error) {
	query := "INSERT INTO carts(user_id, product_id, quantity) VALUES($1, $2, $3) RETURNING *"
	row := c.store.QueryRow(query, newCart.UserID, newCart.ProductID, newCart.Quantity)
	cart, err := scanNewCartRow(row)
	if err != nil {
		log.Printf("Failed to create new cart due to %s", err)
		return nil, utils.ServerError
	}
	return cart, nil
}

func (c *CartModel) Update(cid uint, updateCart *types.UpdateCartRequest) (*Cart, error) {
	query := "UPDATE carts SET quantity=$1 WHERE id=$2 RETURNING *"
	row := c.store.QueryRow(query, updateCart.Quantity, cid)
	cart, err := scanNewCartRow(row)
	if err == sql.ErrNoRows {
		return nil, utils.NotFound
	}
	if err != nil {
		log.Printf("Failed to update cart %d due to %s", cid, err)
		return nil, utils.ServerError
	}
	return cart, nil
}

func (c *CartModel) Delete(cid uint) error {
	query := "UPDATE carts set deleted_at=$1 where id=$2"
	_, err := c.store.Exec(query, time.Now(), cid)
	if err != nil {
		log.Printf("failed to delete cart item with id %d due to %s", cid, err)
	}
	return nil
}

func scanNewCartRow(row *sql.Row) (*Cart, error) {
	cart := Cart{}
	err := row.Scan(
		&cart.ID,
		&cart.UserID,
		&cart.ProductID,
		&cart.Quantity,
		&cart.CreatedAt,
		&cart.UpdatedAt,
		&cart.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func scanCartRows(rows *sql.Rows) ([]*Cart, error) {
	carts := []*Cart{}
	for rows.Next() {
		cart := Cart{}
		err := rows.Scan(
			&cart.ID,
			&cart.UserID,
			&cart.ProductID,
			&cart.Quantity,
			&cart.CreatedAt,
			&cart.UpdatedAt,
			&cart.DeletedAt,
		)
		if err != nil {
			log.Printf("failed to get all products due to %s", err)
			return nil, utils.ServerError
		}

		carts = append(carts, &cart)
	}
	return carts, nil
}

func NewCartModel(store *sql.DB) *CartModel {
	return &CartModel{
		store: store,
	}
}
