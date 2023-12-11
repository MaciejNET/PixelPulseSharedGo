package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"strings"
	"time"
)

var secretKey = []byte("0052ca579e9e4de599e5c95ad3edcb9a")

func extractTokenFromHeader(authorizationHeader string) string {
	if authorizationHeader == "" {
		return ""
	}

	tokenParts := strings.Split(authorizationHeader, " ")
	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
		return ""
	}

	return tokenParts[1]
}

func GenerateToken(user User) (string, error) {
	claims := customClaims{
		ID:               user.ID,
		Email:            user.Email,
		SubscriptionType: user.SubscriptionType,
		Role:             user.Role,
		StandardClaims: jwt.StandardClaims{
			Audience:  "pixel-pulse",
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Id:        uuid.New().String(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "pixel-pulse",
			NotBefore: time.Now().Unix(),
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}
