package api

import (
	"akshidas/e-com/pkg/db"
	"akshidas/e-com/pkg/model"
	"akshidas/e-com/pkg/services"
	"akshidas/e-com/pkg/types"
	"akshidas/e-com/pkg/utils"
	"context"
	"net/http"
)

type ProductCateogriesServicer interface {
	Create(*types.NewProductCategoryRequest) (*model.ProductCategory, error)
	GetAll() ([]*model.ProductCategory, error)
	GetOne(int) (*model.ProductCategory, error)
	Update(int, *types.UpdateProductCategoryRequest) (*model.ProductCategory, error)
	Delete(int) error
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

func (s *ProductCategoriesApi) GetOne(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := parseId(r.PathValue("id"))
	if err != nil {
		return invalidId(w)
	}
	productCategories, err := s.service.GetOne(id)
	if err != nil {
		if err == utils.NotFound {
			return notFound(w)
		}
		return serverError(w)
	}
	return writeJson(w, http.StatusOK, productCategories)
}

func (s *ProductCategoriesApi) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := parseId(r.PathValue("id"))
	if err != nil {
		return invalidId(w)
	}
	updateProductCategory := types.UpdateProductCategoryRequest{}
	if err := DecodeBody(r.Body, &updateProductCategory); err != nil {
		if err == utils.InvalidRequest {
			return invalidRequest(w)
		}
		return serverError(w)
	}

	updatedProductCategory, err := s.service.Update(id, &updateProductCategory)
	if err != nil {
		if err == utils.NotFound {
			return notFound(w)
		}
		return serverError(w)
	}
	return writeJson(w, http.StatusOK, updatedProductCategory)
}

func (s *ProductCategoriesApi) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := parseId(r.PathValue("id"))
	if err != nil {
		return invalidId(w)
	}

	if err := s.service.Delete(id); err != nil {
		if err == utils.NotFound {
			return notFound(w)
		}
		return serverError(w)
	}
	return writeJson(w, http.StatusOK, "delete successfully")
}

func NewProductCategoriesApi(storage *db.Storage) *ProductCategoriesApi {
	model := model.NewProductCategories(storage.DB)
	service := services.NewProductCategoryService(model)
	return &ProductCategoriesApi{service: service}
}