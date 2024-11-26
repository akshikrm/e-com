package services

import (
	"akshidas/e-com/pkg/model"
	"akshidas/e-com/pkg/types"
)

type ProductModeler interface {
	GetAll() ([]*types.ProductsList, error)
	GetOne(int) (*model.Product, error)
	Create(*types.CreateNewProduct) (*model.Product, error)
	Update(int, *types.CreateNewProduct) (*model.Product, error)
	Delete(int) error
}

type ProductService struct {
	productModel ProductModeler
}

func (r *ProductService) Get() ([]*types.ProductsList, error) {
	return r.productModel.GetAll()
}

func (r *ProductService) Create(newProduct *types.CreateNewProduct) error {
	_, err := r.productModel.Create(newProduct)
	return err
}

func (r *ProductService) Update(id int, newProduct *types.CreateNewProduct) (*model.Product, error) {
	return r.productModel.Update(id, newProduct)
}

func (r *ProductService) GetOne(id int) (*model.Product, error) {
	return r.productModel.GetOne(id)
}

func (r *ProductService) Delete(id int) error {
	return r.productModel.Delete(id)
}

func NewProductService(productModel ProductModeler) *ProductService {
	return &ProductService{
		productModel: productModel,
	}
}
