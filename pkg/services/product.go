package services

import (
	"akshidas/e-com/pkg/types"
	"net/url"
)

type ProductStorager interface {
	GetAll(url.Values) ([]*types.ProductsList, error)
	GetOne(int) (*types.Product, error)
	Create(*types.CreateNewProduct) (*types.Product, error)
	Update(int, *types.CreateNewProduct) (*types.Product, error)
	Delete(int) error
}

type ProductService struct {
	productModel ProductStorager
}

func (r *ProductService) Get(filter url.Values) ([]*types.ProductsList, error) {
	return r.productModel.GetAll(filter)
}

func (r *ProductService) Create(newProduct *types.CreateNewProduct) error {
	_, err := r.productModel.Create(newProduct)
	return err
}

func (r *ProductService) Update(id int, newProduct *types.CreateNewProduct) (*types.Product, error) {
	return r.productModel.Update(id, newProduct)
}

func (r *ProductService) GetOne(id int) (*types.Product, error) {
	return r.productModel.GetOne(id)
}

func (r *ProductService) Delete(id int) error {
	return r.productModel.Delete(id)
}

func NewProductService(productModel ProductStorager) *ProductService {
	return &ProductService{
		productModel: productModel,
	}
}
