package api

import (
	"akshidas/e-com/pkg/db"
	"akshidas/e-com/pkg/model"
	"akshidas/e-com/pkg/services"
	"akshidas/e-com/pkg/types"
	"context"
	"fmt"
	"net/http"
)

type ProductServicer interface {
	Get() ([]*model.Product, error)
	GetOne(int) (*model.Product, error)
	Create(*types.CreateNewProduct) error
	Update(int, *types.CreateNewProduct) (*model.Product, error)
	Delete(int) error
}

type ProductApi struct {
	ProductService ProductServicer
}

func (u *ProductApi) GetAll(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	users, err := u.ProductService.Get()
	if err != nil {
		return err
	}
	return writeJson(w, http.StatusOK, users)
}

func (u *ProductApi) GetOne(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := parseId(r.PathValue("id"))
	if err != nil {
		return err
	}
	foundProduct, err := u.ProductService.GetOne(id)
	if err != nil {
		return err
	}
	return writeJson(w, http.StatusOK, foundProduct)
}

func (u *ProductApi) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := parseId(r.PathValue("id"))
	if err != nil {
		return err
	}
	if err := u.ProductService.Delete(id); err != nil {
		return err
	}
	return writeJson(w, http.StatusOK, "deleted successfully")
}

func (u *ProductApi) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	a := &types.CreateNewProduct{}
	if err := DecodeBody(r.Body, &a); err != nil {
		return err
	}
	err := u.ProductService.Create(a)
	if err != nil {
		return err
	}
	return writeJson(w, http.StatusCreated, "product created")
}

func (u *ProductApi) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	a := types.CreateNewProduct{}
	if err := DecodeBody(r.Body, &a); err != nil {
		fmt.Println(err)
		return err
	}
	id, err := parseId(r.PathValue("id"))
	if err != nil {
		return err
	}
	product, err := u.ProductService.Update(id, &a)
	if err != nil {
		return err
	}
	return writeJson(w, http.StatusCreated, product)
}

func NewProductApi(database *db.Storage) *ProductApi {
	productModel := model.NewProductModel(database.DB)
	productService := services.NewProductService(productModel)
	return &ProductApi{ProductService: productService}
}
