package utils

import (
	"errors"
	"time"

	configpkg "websocket-chat/internal/config"

	jwt "github.com/golang-jwt/jwt/v5"
)

// GenerateToken creates a signed JWT for the provided subject and email with the given ttl.
// subject is typically the user ID.
func GenerateToken(subject string, email string, ttl time.Duration, jwtCfg *configpkg.JWTConfig) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":   subject,
		"email": email,
		"iss":   jwtCfg.Issuer,
		"iat":   now.Unix(),
		"exp":   now.Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(jwtCfg.Secret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

// ValidateToken verifies signature and expiration; returns claims if valid.
func ValidateToken(tokenString string, jwtCfg *configpkg.JWTConfig) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtCfg.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims type")
	}
	return claims, nil
}
