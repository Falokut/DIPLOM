package helpers

import (
	"dish_as_a_service/domain"
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(r *http.Request) (string, error) {
	token := r.Header.Get(domain.AuthHeaderName)
	tokenParts := strings.Split(token, " ")
	if len(tokenParts) != 2 ||
		tokenParts[0] != domain.BearerToken ||
		tokenParts[1] == "" {
		return "", domain.DomainInvalidTokenError(errors.New("invalid token")) // nolint:err113,wrapcheck
	}
	return tokenParts[1], nil
}
