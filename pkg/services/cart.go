package services

import (
	"akshidas/e-com/pkg/model"
	"akshidas/e-com/pkg/types"
)

type CartModeler interface {
	GetAll(uint) ([]*model.Cart, error)
	GetOne(uint) (*model.Cart, error)
	Create(*types.CreateCartRequest) (*model.Cart, error)
	Update(uint, *types.UpdateCartRequest) (*model.Cart, error)
	Delete(uint) error
}

type CartService struct {
	cartModel CartModeler
}

func (c *CartService) GetAll(userID uint) ([]*model.Cart, error) {
	return c.cartModel.GetAll(userID)
}

func (c *CartService) GetOne(cid uint) (*model.Cart, error) {
	return c.cartModel.GetOne(cid)
}

func (c *CartService) Create(newCart *types.CreateCartRequest) error {
	_, err := c.cartModel.Create(newCart)
	return err
}

func (c *CartService) Update(cid uint, updateCart *types.UpdateCartRequest) (*model.Cart, error) {
	return c.cartModel.Update(cid, updateCart)
}

func (c *CartService) Delete(cid uint) error {
	return c.cartModel.Delete(cid)
}

func NewCartService(cartModel CartModeler) *CartService {
	return &CartService{cartModel: cartModel}
}
