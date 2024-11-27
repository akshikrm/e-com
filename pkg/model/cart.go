package model

import (
	"akshidas/e-com/pkg/types"
	"akshidas/e-com/pkg/utils"
	"database/sql"
	"log"
	"time"
)

type CartStorage struct {
	store *sql.DB
}

func (c *CartStorage) GetAll(userID uint) ([]*types.Cart, error) {
	query := "SELECT * FROM carts WHERE user_id=$1 AND deleted_at IS NULL"
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

func (c *CartStorage) GetOne(cid uint) (*types.Cart, error) {
	query := "SELECT * FROM carts WHERE id=$1 AND deleted_at IS NULL"
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

func (c *CartStorage) Create(newCart *types.CreateCartRequest) (*types.Cart, error) {
	query := "INSERT INTO carts(user_id, product_id, quantity) VALUES($1, $2, $3) RETURNING *"
	row := c.store.QueryRow(query, newCart.UserID, newCart.ProductID, newCart.Quantity)
	cart, err := scanNewCartRow(row)
	if err != nil {
		log.Printf("Failed to create new cart due to %s", err)
		return nil, utils.ServerError
	}
	return cart, nil
}

func (c *CartStorage) Update(cid uint, updateCart *types.UpdateCartRequest) (*types.Cart, error) {
	query := "UPDATE carts SET quantity=$1 WHERE id=$2 AND deleted_at IS NULL RETURNING *"
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

func (c *CartStorage) Delete(cid uint) error {
	query := "UPDATE carts set deleted_at=$1 where id=$2 AND deleted_at IS NULL"
	if _, err := c.store.Exec(query, time.Now(), cid); err != nil {
		log.Printf("failed to delete cart item with id %d due to %s", cid, err)
		return utils.ServerError
	}
	return nil
}

func scanNewCartRow(row *sql.Row) (*types.Cart, error) {
	cart := types.Cart{}
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

func scanCartRows(rows *sql.Rows) ([]*types.Cart, error) {
	carts := []*types.Cart{}
	for rows.Next() {
		cart := types.Cart{}
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

func NewCartStorage(store *sql.DB) *CartStorage {
	return &CartStorage{
		store: store,
	}
}
