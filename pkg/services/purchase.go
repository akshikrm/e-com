package services

import (
	"akshidas/e-com/pkg/types"
)

type PurchaseStorager interface {
	GetByUserID(uint32) ([]*types.Purchase, error)
	Create([]*types.PurchaseRequest, uint32) error
}

type CartServicer interface {
	GetTotalPriceOfUser(uint) (uint32, error)
	GetAllProductIDByUserID(uint) ([]*uint32, error)
	HardDeleteByUserID(uint) error
}

type PurchaseService struct {
	purchaseStorage PurchaseStorager
	cartService     CartServicer
}

func (s *PurchaseService) Create(newPurchase *types.PurchaseRequest) error {
	cartPrice, err := s.cartService.GetTotalPriceOfUser(newPurchase.UserID)
	if err != nil {
		return err
	}
	productIds, err := s.cartService.GetAllProductIDByUserID(newPurchase.UserID)
	if err != nil {
		return err
	}
	newPurchaseEntry := []*types.PurchaseRequest{}
	for _, productId := range productIds {
		newPurchaseEntry = append(newPurchaseEntry, &types.PurchaseRequest{
			UserID:    newPurchase.UserID,
			ProductID: uint(*productId),
		})
	}
	if err := s.purchaseStorage.Create(newPurchaseEntry, cartPrice); err != nil {
		return err
	}

	return s.cartService.HardDeleteByUserID(newPurchase.UserID)
}

func (s *PurchaseService) GetByUserID(userID uint32) ([]*types.Purchase, error) {
	return s.purchaseStorage.GetByUserID(userID)
}

func NewPurchaseService(purchaseStorage PurchaseStorager, cartService CartServicer) *PurchaseService {
	return &PurchaseService{purchaseStorage: purchaseStorage, cartService: cartService}
}
