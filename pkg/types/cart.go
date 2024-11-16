package types

type CreateCartRequest struct {
	UserID    uint `json:"user_id"`
	ProductID uint `json:"product_id"`
	Quantity  uint `json:"quantity"`
}

type UpdateCartRequest struct {
	Quantity uint `json:"quantity"`
}
