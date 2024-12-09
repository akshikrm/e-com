package types

import "time"

type CreateCartRequest struct {
	UserID    uint `json:"user_id"`
	ProductID uint `json:"product_id"`
	Quantity  uint `json:"quantity"`
}

type UpdateCartRequest struct {
	Quantity uint `json:"quantity"`
}

type Cart struct {
	ID        uint       `json:"id"`
	UserID    uint       `json:"user_id"`
	ProductID uint       `json:"product_id"`
	Quantity  uint       `json:"quantity"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type CartList struct {
	ID       uint `json:"id"`
	Quantity uint `json:"quantity"`
	Product  struct {
		ID          uint   `json:"id"`
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		Price       uint   `json:"price"`
		Description string `json:"description"`
		Image       string `json:"image"`
	} `json:"product"`
	CreatedAt time.Time `json:"created_at"`
}
