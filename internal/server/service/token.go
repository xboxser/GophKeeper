package service

import (
	"fmt"
	"gophkeeper/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Claims — структура утверждений, которая включает стандартные утверждения и
// одно пользовательское UserID
type Claims struct {
	jwt.RegisteredClaims
	User string
}

type TokenService interface {
	BuildJWTString(model.Session) (string, error)
	GetUser(tokenString string) (string, error)
}

type tokenService struct {
	JWTSecret []byte
	TokenExp  time.Duration
}

func NewTokenService(jwtSecret string, tokenExp time.Duration) *tokenService {
	return &tokenService{
		JWTSecret: []byte(jwtSecret),
		TokenExp:  tokenExp,
	}
}

// BuildJWTString - создаёт токен и возвращает его в виде строки.
func (t *tokenService) BuildJWTString(session model.Session) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.TokenExp)),
		},
		User: session.UUID,
	})

	// создаём строку токена
	tokenString, err := token.SignedString(t.JWTSecret)
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}

// GetUserID - возвращает ID пользователя из токена.
func (t *tokenService) GetUser(tokenString string) (string, error) {
	claims := &Claims{}
	jwtSecret := t.JWTSecret
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return jwtSecret, nil
		})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", model.ErrTokenValid
	}

	return claims.User, nil
}
