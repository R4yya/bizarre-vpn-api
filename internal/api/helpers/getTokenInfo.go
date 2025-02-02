package helpers

import (
	"bizarre-vpn-api/internal/lib/jwt"
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetTokenInfo(c *gin.Context) (*jwt.TokenBody, error) {
	op := "internal.api.helpers.GetTokenInfo"

	tokenInfo, isExist := c.Get("tokenInfo")

	if !isExist {
		return nil, fmt.Errorf("%v: tokenInfo not found in gin context", op)
	}

	tokenBody, ok := tokenInfo.(*jwt.TokenBody)
	if !ok {
		return nil, fmt.Errorf("%v: tokenInfo is NOT of type *jwt.TokenBody", op)
	}

	return tokenBody, nil
}
