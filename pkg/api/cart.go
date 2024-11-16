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

type CartServicer interface {
	GetAll(uint) ([]*model.Cart, error)
	GetOne(uint) (*model.Cart, error)
	Create(*types.CreateCartRequest) error
	Update(uint, *types.UpdateCartRequest) (*model.Cart, error)
	Delete(uint) error
}

type CartApi struct {
	cartService CartServicer
}

func (c *CartApi) GetAll(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	userID := ctx.Value("userID")
	carts, err := c.cartService.GetAll(uint(userID.(int)))
	if err != nil {
		return err
	}
	return writeJson(w, http.StatusOK, carts)
}

func (c *CartApi) GetOne(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	cid, err := parseId(r.PathValue("id"))
	if err != nil {
		return invalidId(w)
	}
	cart, err := c.cartService.GetOne(uint(cid))
	if err != nil {
		if err == utils.NotFound {
			return notFound(w)
		}
		return err
	}
	return writeJson(w, http.StatusOK, cart)
}

func (c *CartApi) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	newCart := types.CreateCartRequest{}
	if err := json.NewDecoder(r.Body).Decode(&newCart); err != nil {
		if err == io.EOF {
			return invalidRequest(w)
		}
		return err
	}
	userID := ctx.Value("userID")
	newCart.UserID = uint(userID.(int))
	if err := c.cartService.Create(&newCart); err != nil {
		return err
	}
	return writeJson(w, http.StatusCreated, "cart created")
}

func (c *CartApi) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	cid, err := parseId(r.PathValue("id"))
	if err != nil {
		return invalidId(w)
	}
	updatedCart := types.UpdateCartRequest{}
	if err := json.NewDecoder(r.Body).Decode(&updatedCart); err != nil {
		if err == io.EOF {
			return invalidRequest(w)
		}
		if err == utils.NotFound {
			return notFound(w)
		}
		return err
	}
	cart, err := c.cartService.Update(uint(cid), &updatedCart)
	if err == utils.NotFound {
		return notFound(w)
	}
	if err != nil {
		return err
	}
	return writeJson(w, http.StatusOK, cart)
}

func (c *CartApi) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := parseId(r.PathValue("id"))
	if err != nil {
		return invalidId(w)
	}
	if err := c.cartService.Delete(uint(id)); err != nil {
		if err == utils.NotFound {
			return notFound(w)
		}
		return err
	}
	return writeJson(w, http.StatusOK, "deleted successfully")
}

func NewCartApi(database *db.Storage) *CartApi {
	cartModel := model.NewCartModel(database.DB)
	cartService := services.NewCartService(cartModel)
	return &CartApi{cartService: cartService}
}
