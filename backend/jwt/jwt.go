package jwt

import (
	"dish_as_a_service/domain"
	"dish_as_a_service/entity"
	"time"

	"github.com/pkg/errors"

	jwt2 "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt2.RegisteredClaims
	Value entity.TokenUserInfo `json:"value"`
}

func ParseToken(tokenString string, secret string) (*entity.TokenUserInfo, error) {
	t, err := jwt2.ParseWithClaims(tokenString, &Claims{},
		func(t *jwt2.Token) (any, error) {
			if _, ok := t.Method.(*jwt2.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return []byte(secret), nil
		})

	if err != nil {
		return nil, errors.WithMessage(err, "parse with claims")
	}

	claims, ok := t.Claims.(*Claims)
	if !ok {
		return nil, errors.New("token claims are not of type")
	}

	return &claims.Value, nil
}

func GenerateToken(secret string, tokenTTL time.Duration, value entity.TokenUserInfo) (*domain.TokenResponse, error) {
	expiresAt := time.Now().Add(tokenTTL)
	registeredClaims := jwt2.RegisteredClaims{
		ExpiresAt: &jwt2.NumericDate{Time: expiresAt},
		IssuedAt:  &jwt2.NumericDate{Time: time.Now()}}

	token := jwt2.NewWithClaims(jwt2.SigningMethodHS256, &Claims{registeredClaims, value})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, errors.WithMessage(err, "can't create token")
	}
	return &domain.TokenResponse{
		Token:     tokenString,
		ExpiresAt: expiresAt,
	}, nil
}
