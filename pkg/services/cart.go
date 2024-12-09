package services

import (
	"akshidas/e-com/pkg/types"
)

type CartModeler interface {
	GetAll(uint) ([]*types.CartList, error)
	GetOne(uint) (*types.Cart, error)
	Create(*types.CreateCartRequest) (*types.Cart, error)
	Update(uint, *types.UpdateCartRequest) (*types.Cart, error)
	Delete(uint) error
}

type CartService struct {
	cartModel CartModeler
}

func (c *CartService) GetAll(userID uint) ([]*types.CartList, error) {
	return c.cartModel.GetAll(userID)
}

func (c *CartService) GetOne(cid uint) (*types.Cart, error) {
	return c.cartModel.GetOne(cid)
}

func (c *CartService) Create(newCart *types.CreateCartRequest) error {
	_, err := c.cartModel.Create(newCart)
	return err
}

func (c *CartService) Update(cid uint, updateCart *types.UpdateCartRequest) (*types.Cart, error) {
	return c.cartModel.Update(cid, updateCart)
}

func (c *CartService) Delete(cid uint) error {
	return c.cartModel.Delete(cid)
}

func NewCartService(cartModel CartModeler) *CartService {
	return &CartService{cartModel: cartModel}
}
