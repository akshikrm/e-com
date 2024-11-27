package api

import (
	"akshidas/e-com/pkg/db"
	"akshidas/e-com/pkg/services"
	"akshidas/e-com/pkg/storage"
	"akshidas/e-com/pkg/types"
	"context"
	"net/http"
)

type ProductCateogriesServicer interface {
	Create(*types.NewProductCategoryRequest) (*types.ProductCategory, error)
	GetAll() ([]*types.ProductCategory, error)
	GetNames() ([]*types.ProductCategoryName, error)
	GetOne(int) (*types.ProductCategory, error)
	Update(int, *types.UpdateProductCategoryRequest) (*types.ProductCategory, error)
	Delete(int) error
}

type ProductCategoriesApi struct {
	service ProductCateogriesServicer
}

func (s *ProductCategoriesApi) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	newProductCategory := types.NewProductCategoryRequest{}
	if err := DecodeBody(r.Body, &newProductCategory); err != nil {
		return err
	}
	_, err := s.service.Create(&newProductCategory)
	if err != nil {
		return err
	}
	return writeJson(w, http.StatusCreated, "product category created")
}

func (s *ProductCategoriesApi) GetAll(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	filter := r.URL.Query()
	filterType := filter.Get("type")
	if filterType == "name" {
		productCategoryNames, err := s.service.GetNames()
		if err != nil {
			return err
		}
		return writeJson(w, http.StatusOK, productCategoryNames)
	}
	productCategories, err := s.service.GetAll()
	if err != nil {
		return err
	}
	return writeJson(w, http.StatusOK, productCategories)
}

func (s *ProductCategoriesApi) GetOne(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := parseId(r.PathValue("id"))
	if err != nil {
		return err
	}
	productCategories, err := s.service.GetOne(id)
	if err != nil {
		return err
	}
	return writeJson(w, http.StatusOK, productCategories)
}

func (s *ProductCategoriesApi) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := parseId(r.PathValue("id"))
	if err != nil {
		return err
	}
	updateProductCategory := types.UpdateProductCategoryRequest{}
	if err := DecodeBody(r.Body, &updateProductCategory); err != nil {
		return err
	}

	updatedProductCategory, err := s.service.Update(id, &updateProductCategory)
	if err != nil {
		return err
	}
	return writeJson(w, http.StatusOK, updatedProductCategory)
}

func (s *ProductCategoriesApi) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := parseId(r.PathValue("id"))
	if err != nil {
		return err
	}

	if err := s.service.Delete(id); err != nil {
		return err
	}
	return writeJson(w, http.StatusOK, "delete successfully")
}

func NewProductCategoriesApi(store *db.Storage) *ProductCategoriesApi {
	model := storage.NewProductCategoryStorage(store.DB)
	service := services.NewProductCategoryService(model)
	return &ProductCategoriesApi{service: service}
}
