package auth

import (
	"context"
	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
	"time"
)

var JWKS *keyfunc.JWKS

func InitJWKS(jwksURL string) error {
	var err error
	JWKS, err = keyfunc.Get(jwksURL, keyfunc.Options{})
	return err
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenStr := r.Header.Get("Authorization")
		if !strings.HasPrefix(tokenStr, "Bearer ") {
			http.Error(w, "Missing Bearer token", http.StatusUnauthorized)
			return
		}

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		token, err := jwt.Parse(tokenStr, JWKS.Keyfunc) // main verification happens here
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			role := claims["role"]                                       // adds role to ctx
			ctx := context.WithValue(r.Context(), "user", claims["sub"]) // adds user to ctx
			ctx = context.WithValue(ctx, "role", role)

			exp, ok := claims["exp"].(float64)
			if !ok || int64(exp) < time.Now().Unix() {
				http.Error(w, "Token expired", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Invalid claims", http.StatusUnauthorized)
		}
	})
}
