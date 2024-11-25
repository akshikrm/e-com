package api

import (
	"akshidas/e-com/pkg/db"
	"akshidas/e-com/pkg/model"
	"akshidas/e-com/pkg/services"
	"akshidas/e-com/pkg/types"
	"akshidas/e-com/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"io"
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
		return invalidId(w)
	}
	foundProduct, err := u.ProductService.GetOne(id)
	if err != nil {
		if err == utils.NotFound {
			return notFound(w)
		}
		return err
	}
	return writeJson(w, http.StatusOK, foundProduct)
}

func (u *ProductApi) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := parseId(r.PathValue("id"))
	if err != nil {
		return fmt.Errorf("invalid id")
	}
	if err := u.ProductService.Delete(id); err != nil {
		if err == utils.NotFound {
			return notFound(w)
		}
		return err
	}
	return writeJson(w, http.StatusOK, "deleted successfully")
}

func (u *ProductApi) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	a := &types.CreateNewProduct{}
	if err := json.NewDecoder(r.Body).Decode(a); err != nil {
		if err == io.EOF {
			return invalidRequest(w)
		}
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
	if err := DecodeBody(r.Body, a); err != nil {
		if err == utils.InvalidRequest {
			return invalidRequest(w)
		}
		return err
	}
	id, err := parseId(r.PathValue("id"))
	product, err := u.ProductService.Update(id, &a)
	if err != nil {
		if err == utils.NotFound {
			return notFound(w)
		}
		return err
	}
	return writeJson(w, http.StatusCreated, product)
}

func NewProductApi(database *db.Storage) *ProductApi {
	productModel := model.NewProductModel(database.DB)
	productService := services.NewProductService(productModel)
	return &ProductApi{ProductService: productService}
}
