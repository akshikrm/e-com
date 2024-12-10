package types

import "time"

type PurchaseRequest struct {
	ProductID uint `json:"product_id"`
	UserID    uint `json:"user_id"`
}

type Purchase struct {
	ID        uint32     `json:"id"`
	ProductID uint       `json:"product_id"`
	UserID    uint       `json:"user_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
