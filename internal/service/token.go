package service

import (
	"dekamond-task/internal/config"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type TokenService struct {
	cfg *config.Config
}

func NewTokenService(cfg *config.Config) *TokenService {
	return &TokenService{
		cfg: cfg,
	}
}

func (t *TokenService) GenerateToken(payload map[string]any) (string, error) {
	claims := jwt.MapClaims(payload)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tk, err := token.SignedString([]byte(t.cfg.Token.Secret))
	if err != nil {
		return "", err
	}
	return tk, nil
}
func (t *TokenService) Decode(token string) (map[string]any, error) {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return nil, errors.New("invalid token format")
	}

	payloadSegment := parts[1]

	payloadBytes, err := base64.RawURLEncoding.DecodeString(payloadSegment)
	if err != nil {
		return nil, err
	}

	var payload map[string]any
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, err
	}

	return payload, nil
}

func (t *TokenService) Verify(tk string) (bool, error) {
	token, err := jwt.Parse(tk, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return t.cfg.Token.Secret, nil
	})
	if err != nil {
		return false, err
	}

	return token.Valid, nil
}
