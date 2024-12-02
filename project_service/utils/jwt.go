package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("your_secret_key") // Zameni sa sigurnijim tajnim ključem

type Claims struct {
	UserID  string `json:"userId"` // Dodaj userId u claimove
	Subject string `json:"sub"`
	Role    string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT generiše novi JWT token sa prilagođenim claimovima, uključujući userId
func GenerateJWT(userID string, email string, role string) (string, error) {
	claims := &Claims{
		UserID:  userID, // Postavljamo userId
		Subject: email,
		Role:    role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token ističe za 24 sata
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateJWT validira JWT token i vraća userId, subject (email) i role ako je validan
func ValidateJWT(tokenStr string) (string, string, string, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return "", "", "", err
	}

	return claims.UserID, claims.Subject, claims.Role, nil
}
