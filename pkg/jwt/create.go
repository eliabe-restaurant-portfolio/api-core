package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type AccessDetails struct {
	Token     string `json:"access_token"`
	ExpiresAt string `json:"expires_at"`
	IssuedAt  string `json:"issued_at"`
}

type JwtCreateInput struct {
	Content        string
	Duration       time.Duration
	AccessClientId string
}

func Create(input JwtCreateInput) (*AccessDetails, error) {
	if input.Duration <= 0 {
		return nil, fmt.Errorf("invalid TTL: duration must be positive")
	}
	if input.AccessClientId == "" {
		return nil, fmt.Errorf("access client ID is empty")
	}

	// Load private key from file or .env
	var privatePem []byte
	privatePem, err := os.ReadFile("./storage/private_key.pem")
	if err != nil {
		// Fallback to AUTH_PRIVATE_KEY from .env
		if envKey := os.Getenv("AUTH_PRIVATE_KEY"); envKey != "" {
			privatePem = []byte(envKey)
		} else {
			return nil, fmt.Errorf("failed to read private key file ./storage/private_key.pem: %w", err)
		}
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePem)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA private key: %w", err)
	}

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(input.Duration)

	claims := jwt.MapClaims{
		"aud":    input.AccessClientId,
		"jti":    uuid.NewString(),
		"iat":    jwt.NewNumericDate(issuedAt),
		"nbf":    jwt.NewNumericDate(issuedAt),
		"exp":    jwt.NewNumericDate(expiresAt),
		"sub":    input.Content,
		"scopes": []string{},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	accessToken, err := token.SignedString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign JWT: %w", err)
	}

	return &AccessDetails{
		Token:     accessToken,
		IssuedAt:  issuedAt.Format(time.RFC3339),
		ExpiresAt: expiresAt.Format(time.RFC3339),
	}, nil
}
