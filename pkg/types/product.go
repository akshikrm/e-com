package types

type CreateNewProduct struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Price       uint   `json:"price"`
	Image       string `json:"image"`
	Description string `json:"description"`
}
