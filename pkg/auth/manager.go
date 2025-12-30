package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenManager interface {
	NewJWT(userID int64, ttl time.Duration) (string, error)
	Parse(accessToken string) (int64, error)
}

type Manager struct {
	signingKey string
}

func NewManager(signingKey string) (*Manager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}
	return &Manager{signingKey: signingKey}, nil
}

func (m *Manager) NewJWT(userID int64, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(ttl).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(m.signingKey))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (m *Manager) Parse(accessToken string) (int64, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.signingKey), nil
	})
	if err != nil {
		return 0, err
	}
	if !token.Valid {
		return 0, fmt.Errorf("token is invalid")
	}
	var sub float64
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("smth went wrong")
	}
	sub = claims["sub"].(float64)
	return int64(sub), nil
}
