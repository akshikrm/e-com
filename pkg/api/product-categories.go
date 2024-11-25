package api

import (
	"akshidas/e-com/pkg/db"
	"akshidas/e-com/pkg/model"
	"akshidas/e-com/pkg/services"
	"akshidas/e-com/pkg/types"
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
		return err
	}
	_, err := s.service.Create(&newProductCategory)
	if err != nil {
		return err
	}
	return writeJson(w, http.StatusCreated, "product category created")
}

func (s *ProductCategoriesApi) GetAll(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
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

func NewProductCategoriesApi(storage *db.Storage) *ProductCategoriesApi {
	model := model.NewProductCategories(storage.DB)
	service := services.NewProductCategoryService(model)
	return &ProductCategoriesApi{service: service}
}
