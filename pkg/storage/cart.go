package storage

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

func (c *CartStorage) GetAll(userID uint) ([]*types.CartList, error) {
	query := "SELECT c.id, c.quantity, p.id, p.name, p.slug, p.price, p.description, p.image, c.created_at FROM carts c INNER JOIN products p ON c.product_id=p.id WHERE c.user_id=$1 AND c.deleted_at IS NULL"
	rows, err := c.store.Query(query, userID)
	if err == sql.ErrNoRows {
		return nil, utils.NotFound
	}
	if err != nil {
		log.Printf("failed to get all carts due to %s", err)
		return nil, utils.ServerError
	}
	carts := []*types.CartList{}
	for rows.Next() {
		cart := types.CartList{}
		err := rows.Scan(
			&cart.ID,
			&cart.Quantity,
			&cart.Product.ID,
			&cart.Product.Name,
			&cart.Product.Slug,
			&cart.Product.Price,
			&cart.Product.Description,
			&cart.Product.Image,
			&cart.CreatedAt,
		)
		if err != nil {
			log.Printf("failed to scan carts due to %s", err)
			return nil, utils.ServerError
		}
		carts = append(carts, &cart)
	}
	return carts, nil
}

func (c *CartStorage) GetOne(cid uint) (*types.CartList, error) {
	query := "SELECT c.id, c.quantity, p.id, p.name, p.slug, p.price, p.description, p.image, c.created_at FROM carts c INNER JOIN products p ON c.product_id=p.id WHERE c.id=$1 AND c.deleted_at IS NULL"
	row := c.store.QueryRow(query, cid)
	cart := types.CartList{}
	err := row.Scan(
		&cart.ID,
		&cart.Quantity,
		&cart.Product.ID,
		&cart.Product.Name,
		&cart.Product.Slug,
		&cart.Product.Price,
		&cart.Product.Description,
		&cart.Product.Image,
		&cart.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, utils.NotFound
	}

	if err != nil {
		log.Printf("failed to scan carts due to %s", err)
		return nil, utils.ServerError
	}
	return &cart, nil
}

func (c *CartStorage) CheckIfEntryExist(userID, productID uint) (bool, error) {
	query := "select exists(select 1 from carts where user_id=$1 and product_id=$2 and deleted_at IS NULL)"
	row := c.store.QueryRow(query, userID, productID)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		log.Printf("failed to scan due to %s", err)
		return false, utils.ServerError
	}
	return exists, nil
}

func (c *CartStorage) UpdateQuantity(cid, pid, qty uint) error {
	query := "UPDATE carts SET quantity=quantity+$1 WHERE user_id=$2 and product_id=$3"
	if _, err := c.store.Exec(query, qty, cid, pid); err != nil {
		if err == sql.ErrNoRows {
			return utils.NotFound
		}
		log.Printf("Failed to update cart %d due to %s", cid, err)
		return utils.ServerError
	}

	return nil
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

func (c *CartStorage) Update(cid uint, updateCart *types.UpdateCartRequest) (*types.CartList, error) {
	query := "UPDATE carts SET quantity=$1 WHERE id=$2 AND deleted_at IS NULL"
	if _, err := c.store.Exec(query, updateCart.Quantity, cid); err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NotFound
		}
		log.Printf("Failed to update cart %d due to %s", cid, err)
		return nil, utils.ServerError
	}

	return c.GetOne(cid)
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
