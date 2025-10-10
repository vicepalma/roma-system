package security

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenPair struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

func mustEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func GenerateTokens(userID string) (TokenPair, error) {
	secret := os.Getenv("JWT_SECRET")
	accessTTL := time.Duration(mustEnvInt("ACCESS_TTL_MIN", 15)) * time.Minute
	refreshTTL := time.Duration(mustEnvInt("REFRESH_TTL_H", 168)) * time.Hour

	now := time.Now()

	access := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"typ": "access",
		"iat": now.Unix(),
		"exp": now.Add(accessTTL).Unix(),
	})
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"typ": "refresh",
		"iat": now.Unix(),
		"exp": now.Add(refreshTTL).Unix(),
	})

	accessStr, err := access.SignedString([]byte(secret))
	if err != nil {
		return TokenPair{}, err
	}
	refreshStr, err := refresh.SignedString([]byte(secret))
	if err != nil {
		return TokenPair{}, err
	}
	return TokenPair{AccessToken: accessStr, RefreshToken: refreshStr}, nil
}

func ParseAndValidate(tokenStr string) (*jwt.Token, jwt.MapClaims, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secret, nil
	})
	if err != nil {
		return nil, nil, err
	}
	claims, _ := tok.Claims.(jwt.MapClaims)
	return tok, claims, nil
}
