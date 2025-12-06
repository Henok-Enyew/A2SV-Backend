package infrastructure

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTGenerator struct {
	secretKey string
}

func NewJWTGenerator() *JWTGenerator {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "your-secret-key-change-in-production"
	}
	return &JWTGenerator{secretKey: secretKey}
}

func (j *JWTGenerator) Generate(userID, username, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTGenerator) Validate(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		result := make(map[string]interface{})
		result["user_id"] = claims["user_id"]
		result["username"] = claims["username"]
		result["role"] = claims["role"]
		return result, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

