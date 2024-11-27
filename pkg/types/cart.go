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
