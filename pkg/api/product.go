package api

import (
	"akshidas/e-com/pkg/model"
	"akshidas/e-com/pkg/types"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type ProductServicer interface {
	Get() ([]*model.Product, error)
	Create(*types.CreateNewProduct) error
}

type ProductApi struct {
	ProductService ProductServicer
}

func (u *ProductApi) GetAll(w http.ResponseWriter, r *http.Request) error {
	users, err := u.ProductService.Get()
	if err != nil {
		return err
	}
	return writeJson(w, http.StatusOK, users)
}

func (u *ProductApi) Create(w http.ResponseWriter, r *http.Request) error {
	a := &types.CreateNewProduct{}
	if err := json.NewDecoder(r.Body).Decode(a); err != nil {
		if err == io.EOF {
			return errors.New("invalid request")
		}
		return err
	}
	err := u.ProductService.Create(a)
	if err != nil {
		return err
	}
	return writeJson(w, http.StatusCreated, "product created")
}

func NewProductApi(productService ProductServicer) *ProductApi {
	return &ProductApi{ProductService: productService}
}
