package krap

import (
	"github.com/roidaradal/fn/check"
	"github.com/roidaradal/fn/str"
)

const (
	authTokenGlue    string = "/"
	SESSION_ACTIVE   string = "ACTIVE"
	SESSION_EXTENDED string = "EXTENDED"
	SESSION_EXPIRED  string = "EXPIRED"
	SESSION_LOGOUT   string = "LOGOUT"
)

type AuthParams struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

type AuthToken struct {
	Type string `validate:"required"`
	Code string `validate:"required"`
}

func (a AuthToken) String() string {
	return a.Type + authTokenGlue + a.Code
}

func NewAuthToken(authToken string) *AuthToken {
	parts := str.CleanSplit(authToken, authTokenGlue)
	if len(parts) != 2 || check.Any(parts, check.IsBlankString) {
		return nil
	}
	return &AuthToken{
		Type: parts[0],
		Code: parts[1],
	}
}

func IsAuthToken(authToken string) bool {
	parts := str.CleanSplit(authToken, authTokenGlue)
	return len(parts) == 2 && check.All(parts, check.IsNotBlankString)
}
