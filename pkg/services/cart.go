package services

import (
	"akshidas/e-com/pkg/types"
	"fmt"
)

type CartModeler interface {
	GetAll(uint) ([]*types.CartList, error)
	GetOne(uint) (*types.CartList, error)
	Create(*types.CreateCartRequest) (*types.Cart, error)
	Update(uint, *types.UpdateCartRequest) (*types.CartList, error)
	Delete(uint) error
	CheckIfEntryExist(uint, uint) (bool, error)
	UpdateQuantity(uint, uint, uint) error
	GetTotalPriceOfUser(uint) (uint32, error)
	GetAllProductIDByUserID(uint) ([]*uint32, error)
	HardDeleteByUserID(uint) error
}

type CartService struct {
	cartModel CartModeler
}

func (c *CartService) GetAll(userID uint) ([]*types.CartList, error) {
	return c.cartModel.GetAll(userID)
}

func (c *CartService) GetTotalPriceOfUser(userID uint) (uint32, error) {
	return c.cartModel.GetTotalPriceOfUser(userID)
}
func (c *CartService) GetAllProductIDByUserID(userID uint) ([]*uint32, error) {
	return c.cartModel.GetAllProductIDByUserID(userID)
}

func (c *CartService) GetOne(cid uint) (*types.CartList, error) {
	return c.cartModel.GetOne(cid)
}

func (c *CartService) Create(newCart *types.CreateCartRequest) error {
	exists, err := c.cartModel.CheckIfEntryExist(newCart.UserID, newCart.ProductID)
	if err != nil {
		return err
	}

	fmt.Println(exists)
	if exists {
		return c.cartModel.UpdateQuantity(newCart.UserID, newCart.ProductID, newCart.Quantity)
	}

	_, err = c.cartModel.Create(newCart)
	return err
}

func (c *CartService) Update(cid uint, updateCart *types.UpdateCartRequest) (*types.CartList, error) {
	return c.cartModel.Update(cid, updateCart)
}

func (c *CartService) Delete(cid uint) error {
	return c.cartModel.Delete(cid)
}

func (c *CartService) HardDeleteByUserID(userID uint) error {
	return c.cartModel.HardDeleteByUserID(userID)
}

func NewCartService(cartModel CartModeler) *CartService {
	return &CartService{cartModel: cartModel}
}
