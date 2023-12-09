package auth

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := extractTokenFromHeader(r.Header.Get("Authorization"))
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims := &customClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func AdminOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claimsValue := r.Context().Value("claims")
		if claimsValue == nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		claims, ok := claimsValue.(*customClaims)
		if !ok || claims.Role != "admin" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func PremiumOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claimsValue := r.Context().Value("claims")
		if claimsValue == nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		claims, ok := claimsValue.(*customClaims)
		if !ok || claims.SubscriptionType != "premium" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
