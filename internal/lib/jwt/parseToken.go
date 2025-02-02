package jwt

import (
	"fmt"

	jwt "github.com/golang-jwt/jwt/v5"
)

func ParseToken(tokenString string, secretKey []byte) (*TokenBody, error) {
	op := "internal.lib.jwt.parseToken"

	token, err := jwt.ParseWithClaims(tokenString, &ExpandedTokenBody{}, func(token *jwt.Token) (interface{}, error) {
		// since we only use the one private key to sign the tokens,
		// we also only use its public counter part to verify
		return secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("%v: error of parsing token %w", op, err)
	}

	claims := token.Claims.(*ExpandedTokenBody)

	return &claims.TokenBody, nil
}
