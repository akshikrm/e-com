package api

import (
	"akshidas/e-com/pkg/utils"
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type apiFuncWithContext func(context.Context, http.ResponseWriter, *http.Request) error
type AdminWrapperFunc func(context.Context, apiFuncWithContext) apiFunc

func IsAdmin(userService UserServicer) AdminWrapperFunc {
	return func(ctx context.Context, f apiFuncWithContext) apiFunc {
		validateAdmin := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			id := ctx.Value("userID")
			user, err := userService.GetOne(id.(int))
			if err != nil {
				return accessDenied(w)
			}
			if user.Role == "admin" {
				return f(ctx, w, r)
			}
			return accessDenied(w)
		}

		return IsAuthenticated(ctx, validateAdmin)
	}
}

func IsAuthenticated(ctx context.Context, f apiFuncWithContext) apiFunc {
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
			ctx = context.WithValue(ctx, "userID", id)
			return f(ctx, w, r)
		}
		return accessDenied(w)
	}
}
