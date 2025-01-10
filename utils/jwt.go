package utils

import (
	"laundre/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("your-secret-key")

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	BranchID *uint  `json:"branch_id"`
	*jwt.RegisteredClaims
}

func GenerateToken(user models.User) (string, error) {

	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		BranchID: user.BranchID,
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenExpired
	}

	return claims, nil
}
