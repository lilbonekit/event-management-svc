package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getSecretKey() ([]byte, error) {
	key := os.Getenv("JWT_SECRET")
	if key == "" {
		return nil, errors.New("JWT_SECRET not set")
	}
	return []byte(key), nil
}

func GenerateToken(email string, userID int64) (string, error) {
	secretKey, err := getSecretKey()

	if err != nil {
		return "", err
	}

	fmt.Printf("Generating token for userID: %d\n", userID)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 2).Unix(),
	})

	fmt.Printf("Token claims set for userID: %d\n", userID)

	return token.SignedString([]byte(secretKey))
}

func ValidateToken(tokenString string) (int64, error) {
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return getSecretKey()
	})

	if err != nil {
		return 0, err
	}

	isTokenValid := parsedToken.Valid
	if !isTokenValid {
		return 0, jwt.ErrTokenUnverifiable
	}

	// Example of extracting claims if needed
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, jwt.ErrTokenInvalidClaims
	}

	userID := int64(claims["userID"].(float64))

	return userID, nil
}
