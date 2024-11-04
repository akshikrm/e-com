package services

import (
	"akshidas/e-com/pkg/model"
	"akshidas/e-com/pkg/types"
	"database/sql"
)

type ProductModeler interface {
	GetAll() ([]*model.Product, error)
	// GetOne(int) (*model.Product, error)
	Create(*types.CreateNewProduct) (uint, error)
	// Update(int, *types.CreateProductRequest) (*model.Product, error)
	// Delete(int) error
}

type ProductService struct {
	roleModel ProductModeler
}

func (r *ProductService) Get() ([]*model.Product, error) {
	return r.roleModel.GetAll()
}

func (r *ProductService) Create(newProduct *types.CreateNewProduct) error {
	_, err := r.roleModel.Create(newProduct)
	return err
}

// func (r *ProductService) GetOne(id int) (*model.Product, error) {
// 	return r.roleModel.GetOne(id)
// }
//
// func (r *ProductService) Update(id int, newProduct *types.CreateProductRequest) (*model.Product, error) {
// 	return r.roleModel.Update(id, newProduct)
// }
//
// func (r *ProductService) Delete(id int) error {
// 	return r.roleModel.Delete(id)
// }
//

func NewProductService(database *sql.DB) *ProductService {
	roleModel := model.NewProductModel(database)
	return &ProductService{
		roleModel: roleModel,
	}
}
