package api

import (
	"akshidas/e-com/pkg/db"
	"akshidas/e-com/pkg/model"
	"akshidas/e-com/pkg/services"
	"akshidas/e-com/pkg/types"
	"akshidas/e-com/pkg/utils"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type ProductCateogriesServicer interface {
	Create(*types.NewProductCategoryRequest) (*model.ProductCategory, error)
	GetAll() ([]*model.ProductCategory, error)
}

type ProductCategoriesApi struct {
	service ProductCateogriesServicer
}

func (s *ProductCategoriesApi) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	newProductCategory := types.NewProductCategoryRequest{}
	if err := DecodeBody(r.Body, &newProductCategory); err != nil {
		if err == utils.InvalidRequest {
			return invalidRequest(w)
		}
		return serverError(w)
	}
	_, err := s.service.Create(&newProductCategory)
	if err != nil {
		if err == utils.NotFound {
			return notFound(w)
		}
		return serverError(w)
	}
	return writeJson(w, http.StatusCreated, "product category created")
}

func (s *ProductCategoriesApi) GetAll(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	productCategories, err := s.service.GetAll()
	if err != nil {
		if err == utils.NotFound {
			return notFound(w)
		}
		return serverError(w)
	}
	return writeJson(w, http.StatusOK, productCategories)

}

func DecodeBody(body io.ReadCloser, a any) error {
	if err := json.NewDecoder(body).Decode(a); err != nil {
		if err == io.EOF {
			return utils.InvalidRequest
		}
		return err
	}
	return nil
}

func NewProductCategoriesApi(storage *db.Storage) *ProductCategoriesApi {
	model := model.NewProductCategories(storage.DB)

	service := services.NewProductCategoryService(model)
	return &ProductCategoriesApi{service: service}
}
