package auth

import "github.com/dgrijalva/jwt-go"

type customClaims struct {
	ID               string `json:"id"`
	Email            string `json:"email"`
	SubscriptionType string `json:"subscription_type"`
	Role             string `json:"role"`
	jwt.StandardClaims
}
