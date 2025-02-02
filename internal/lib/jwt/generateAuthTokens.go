package jwt

import (
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type TokenBody struct {
	UserID int64  `json:"userId"`
	Role   string `json:"role"`
}

type ExpandedTokenBody struct {
	TokenBody
	jwt.RegisteredClaims
}

func GenerateAuthTokens(userID int64, userRole string, accessSecretKey []byte, refreshSecretKey []byte) (string, string, error) {
	op := "internal.lib.jwt.GenerateAuthTokens"

	accessTokenPrepared := jwt.NewWithClaims(jwt.SigningMethodHS256, ExpandedTokenBody{
		TokenBody{
			userID,
			userRole,
		},
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	accessToken, err := accessTokenPrepared.SignedString(accessSecretKey)

	if err != nil {
		return "", "", fmt.Errorf("%v: error signed access token %w", op, err)
	}

	refreshTokenPrepared := jwt.NewWithClaims(jwt.SigningMethodHS256, ExpandedTokenBody{
		TokenBody{
			userID,
			userRole,
		},
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	refreshToken, err := refreshTokenPrepared.SignedString(refreshSecretKey)

	if err != nil {
		return "", "", fmt.Errorf("%v: error signed refresh token %w", op, err)
	}

	return accessToken, refreshToken, nil
}
