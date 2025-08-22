package jwt

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

type JwtReadInput struct {
	ExternalToken string
	PublicPem     []byte
}

func Read(input JwtReadInput) (jwt.MapClaims, error) {
	if len(input.ExternalToken) == 0 {
		return nil, fmt.Errorf("JWT token is empty")
	}
	if len(input.PublicPem) == 0 {
		return nil, fmt.Errorf("public key is empty")
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(input.PublicPem)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA public key from PEM: %w", err)
	}

	token, err := jwt.Parse(input.ExternalToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse or verify JWT: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("JWT is invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid JWT claims type, expected MapClaims")
	}

	return claims, nil
}

func LoadPublicKeyFromFile(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file %s: %w", filePath, err)
	}
	return data, nil
}
