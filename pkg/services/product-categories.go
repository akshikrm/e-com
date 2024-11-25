package services

import (
	"akshidas/e-com/pkg/model"
	"akshidas/e-com/pkg/types"
)

type ProductCategoriesModeler interface {
	Create(*types.NewProductCategoryRequest) (*model.ProductCategory, error)
	GetAll() ([]*model.ProductCategory, error)
	GetOne(int) (*model.ProductCategory, error)
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

func NewProductCategoryService(model ProductCategoriesModeler) *ProductCategoryService {
	return &ProductCategoryService{
		model: model,
	}
}
