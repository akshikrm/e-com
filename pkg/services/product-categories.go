package services

import (
	"akshidas/e-com/pkg/types"
)

type ProductCategoriesModeler interface {
	Create(*types.NewProductCategoryRequest) (*types.ProductCategory, error)
	GetAll() ([]*types.ProductCategory, error)
	GetOne(int) (*types.ProductCategory, error)
	Update(int, *types.UpdateProductCategoryRequest) (*types.ProductCategory, error)
	Delete(int) error
}

type ProductCategoryService struct {
	model ProductCategoriesModeler
}

func (p *ProductCategoryService) Create(newCategory *types.NewProductCategoryRequest) (*types.ProductCategory, error) {
	return p.model.Create(newCategory)
}

func (p *ProductCategoryService) GetAll() ([]*types.ProductCategory, error) {
	return p.model.GetAll()
}

func (p *ProductCategoryService) GetOne(id int) (*types.ProductCategory, error) {
	return p.model.GetOne(id)
}

func (p *ProductCategoryService) Update(id int, updateProductCategory *types.UpdateProductCategoryRequest) (*types.ProductCategory, error) {
	return p.model.Update(id, updateProductCategory)
}

func (p *ProductCategoryService) Delete(id int) error {
	return p.model.Delete(id)
}

func NewProductCategoryService(model ProductCategoriesModeler) *ProductCategoryService {
	return &ProductCategoryService{
		model: model,
	}
}
