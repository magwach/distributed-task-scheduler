package auth

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateToken(userId, email, role string) (string, error) {

	claims := Claims{
		UserID: userId,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(jwtSecret)

	if err != nil {
		log.Println("Failed to generate jwt: ", err)
		return "", err
	}

	return signedToken, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		log.Println("Failed to validate token: ", err)
		return nil, err
	}

	if !token.Valid {
		log.Println("Invalid token")
		return nil, errors.New("invalid token")
	}

	return claims, nil

}
