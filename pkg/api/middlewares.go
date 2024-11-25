package api

import (
	"akshidas/e-com/pkg/utils"
	"context"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

type apiFuncWithContext func(context.Context, http.ResponseWriter, *http.Request) error

type MiddleWares struct{ userService UserServicer }

func (m *MiddleWares) IsAdmin(ctx context.Context, f apiFuncWithContext) apiFunc {
	validateAdmin := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		if role := ctx.Value("role"); role != "admin" {
			return accessDenied(w)
		}
		return f(ctx, w, r)
	}
	return m.IsAuthenticated(ctx, validateAdmin)
}

func (m *MiddleWares) IsAuthenticated(ctx context.Context, f apiFuncWithContext) apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		authtoken := r.Header.Get("Authorization")
		token, err := utils.ValidateJWT(authtoken)
		if err != nil {
			return accessDenied(w)
		}
		if !token.Valid {
			return accessDenied(w)
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			id := int(claims["sub"].(float64))
			user, err := m.userService.GetOne(id)

			if err != nil {
				return accessDenied(w)
			}
			ctx = context.WithValue(ctx, "userID", id)
			ctx = context.WithValue(ctx, "role", user.Role)
			return f(ctx, w, r)
		}
		return accessDenied(w)
	}
}

func NewMiddleWare(userService UserServicer) *MiddleWares {
	return &MiddleWares{userService: userService}
}
