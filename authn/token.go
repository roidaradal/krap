package authn

import (
	"github.com/gin-gonic/gin"
	"github.com/roidaradal/fn/check"
	"github.com/roidaradal/fn/str"
	"github.com/roidaradal/krap"
	"github.com/roidaradal/rdb/ze"
)

const (
	authTokenGlue string = "/"
	authHeaderKey string = "Authorization"
)

// Creates authn.Token from string "Type/Code"
func NewToken(authToken string) *Token {
	parts := str.CleanSplit(authToken, authTokenGlue)
	if len(parts) != 2 || check.Any(parts, check.IsEmptyString) {
		return nil
	}
	return &Token{
		Type: parts[0],
		Code: parts[1],
	}
}

// Checks if authToken string can be a valid authn.Token
func IsToken(authToken string) bool {
	parts := str.CleanSplit(authToken, authTokenGlue)
	return len(parts) == 2 && check.All(parts, check.NotEmptyString)
}

// Get the authn.Token from the Authorization header
func WebAuthToken(c *gin.Context) *Token {
	authHeader := c.GetHeader(authHeaderKey)
	return NewToken(authHeader)
}

// Get the authn.Token from the Authorizatio header;
// On error, send error response
func ReqAuthToken(c *gin.Context, response *krap.ResponseType) *Token {
	authToken := WebAuthToken(c)
	if authToken == nil {
		rq := &ze.Request{Status: ze.Err401}
		response.SendErrorFn(c, rq, ErrInvalidSession)
		return nil
	}
	return authToken
}
