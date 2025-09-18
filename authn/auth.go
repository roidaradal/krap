package authn

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/roidaradal/fn/check"
	"github.com/roidaradal/fn/str"
	"github.com/roidaradal/krap"
)

const (
	authTokenGlue    string = "/"
	authHeaderKey    string = "Authorization"
	SESSION_ACTIVE   string = "ACTIVE"
	SESSION_EXTENDED string = "EXTENDED"
	SESSION_EXPIRED  string = "EXPIRED"
	SESSION_LOGOUT   string = "LOGOUT"
)

var errInvalidSession = errors.New("public: Invalid session")

type Params struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

type Token struct {
	Type string `validate:"required"`
	Code string `validate:"required"`
}

func (a Token) String() string {
	return a.Type + authTokenGlue + a.Code
}

func NewToken(authToken string) *Token {
	parts := str.CleanSplit(authToken, authTokenGlue)
	if len(parts) != 2 || check.Any(parts, check.IsBlankString) {
		return nil
	}
	return &Token{
		Type: parts[0],
		Code: parts[1],
	}
}

func IsToken(authToken string) bool {
	parts := str.CleanSplit(authToken, authTokenGlue)
	return len(parts) == 2 && check.All(parts, check.IsNotBlankString)
}

func WebAuthToken(c *gin.Context) *Token {
	authHeader := c.GetHeader(authHeaderKey)
	return NewToken(authHeader)
}

func ReqAuthToken(c *gin.Context, resp *krap.ResponseType) *Token {
	authToken := WebAuthToken(c)
	if authToken == nil {
		resp.SendErrorFn(c, nil, errInvalidSession)
		return nil
	}
	return authToken
}
