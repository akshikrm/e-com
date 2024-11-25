package services

import (
	"akshidas/e-com/pkg/model"
	"akshidas/e-com/pkg/types"
)

type ProductCategoriesModeler interface {
	Create(*types.NewProductCategoryRequest) (*model.ProductCategory, error)
	GetAll() ([]*model.ProductCategory, error)
	GetOne(int) (*model.ProductCategory, error)
	Update(int, *types.UpdateProductCategoryRequest) (*model.ProductCategory, error)
	Delete(int) error
}

type ProductCategoryService struct {
	model ProductCategoriesModeler
}

func (p *ProductCategoryService) Create(newCategory *types.NewProductCategoryRequest) (*model.ProductCategory, error) {
	return p.model.Create(newCategory)
}

func (p *ProductCategoryService) GetAll() ([]*model.ProductCategory, error) {
	return p.model.GetAll()
}

func (p *ProductCategoryService) GetOne(id int) (*model.ProductCategory, error) {
	return p.model.GetOne(id)
}

func (p *ProductCategoryService) Update(id int, updateProductCategory *types.UpdateProductCategoryRequest) (*model.ProductCategory, error) {
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
