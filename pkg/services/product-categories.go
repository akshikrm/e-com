package services

import (
	"akshidas/e-com/pkg/types"
)

type ProductCategoriesStorager interface {
	Create(*types.NewProductCategoryRequest) (*types.ProductCategory, error)
	GetNames() ([]*types.ProductCategoryName, error)
	GetAll() ([]*types.ProductCategory, error)
	GetOne(int) (*types.ProductCategory, error)
	Update(int, *types.UpdateProductCategoryRequest) (*types.ProductCategory, error)
	Delete(int) error
}

type ProductCategoryService struct {
	storage ProductCategoriesStorager
}

func (p *ProductCategoryService) Create(newCategory *types.NewProductCategoryRequest) (*types.ProductCategory, error) {
	return p.storage.Create(newCategory)
}

func (p *ProductCategoryService) GetNames() ([]*types.ProductCategoryName, error) {
	return p.storage.GetNames()
}
func (p *ProductCategoryService) GetAll() ([]*types.ProductCategory, error) {
	return p.storage.GetAll()
}

func (p *ProductCategoryService) GetOne(id int) (*types.ProductCategory, error) {
	return p.storage.GetOne(id)
}

func (p *ProductCategoryService) Update(id int, updateProductCategory *types.UpdateProductCategoryRequest) (*types.ProductCategory, error) {
	return p.storage.Update(id, updateProductCategory)
}

func (p *ProductCategoryService) Delete(id int) error {
	return p.storage.Delete(id)
}

func NewProductCategoryService(storage ProductCategoriesStorager) *ProductCategoryService {
	return &ProductCategoryService{
		storage: storage,
	}
}
