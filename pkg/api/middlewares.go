package api

import (
	"akshidas/e-com/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

type apiFuncWithContext func(int, http.ResponseWriter, *http.Request) error
type AdminWrapperFunc func(apiFuncWithContext) apiFunc

func IsAdmin(userService UserServicer) AdminWrapperFunc {
	return func(f apiFuncWithContext) apiFunc {
		validateAdmin := func(id int, w http.ResponseWriter, r *http.Request) error {
			user, err := userService.GetOne(id)
			if err != nil {
				return accessDenied(w)
			}
			if user.Role == "admin" {
				return f(id, w, r)
			}
			return accessDenied(w)
		}

		return IsAuthenticated(validateAdmin)
	}
}

func IsAuthenticated(f apiFuncWithContext) apiFunc {
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
			return f(id, w, r)
		}
		return accessDenied(w)
	}
}
