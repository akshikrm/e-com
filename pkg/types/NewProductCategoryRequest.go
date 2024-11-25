package types

type NewProductCategoryRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}
