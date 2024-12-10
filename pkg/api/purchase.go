package api

import (
	"akshidas/e-com/pkg/db"
	"akshidas/e-com/pkg/services"
	"akshidas/e-com/pkg/storage"
	"akshidas/e-com/pkg/types"
	"context"
	"net/http"
)

type PurchaseServecier interface {
	Create(*types.PurchaseRequest) error
	GetByUserID(uint32) ([]*types.Purchase, error)
}

type PurchaseApi struct {
	service PurchaseServecier
}

func (a *PurchaseApi) GetByUserID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := ctx.Value("userID")
	purchases, err := a.service.GetByUserID(id.(uint32))
	if err != nil {
		return err
	}
	return writeJson(w, http.StatusOK, purchases)
}

func (a *PurchaseApi) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := ctx.Value("userID").(int)
	purchasePayload := types.PurchaseRequest{
		UserID: uint(id),
	}

	err := a.service.Create(&purchasePayload)
	if err != nil {
		return err
	}

	return writeJson(w, http.StatusOK, "order placed")
}

func NewPurchaseApi(database *db.Storage) *PurchaseApi {
	purchaseStorage := storage.NewPurchaseStorage(database.DB)
	cartStorage := storage.NewCartStorage(database.DB)
	purchaseService := services.NewPurchaseService(purchaseStorage, cartStorage)
	return &PurchaseApi{service: purchaseService}
}
