package types

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
