package types

import "time"

type CreateNewProduct struct {
	Name        string `json:"name"`
	CategoryID  uint   `json:"category_id"`
	Slug        string `json:"slug"`
	Price       uint   `json:"price"`
	Image       string `json:"image"`
	Description string `json:"description"`
}

type ProductsList struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Price       uint   `json:"price"`
	Image       string `json:"image"`
	Description string `json:"description"`
	Category    struct {
		ID          uint   `json:"id"`
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
	} `json:"category"`
}

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
