package jwt

import (
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type TokenBody struct {
	UserID int64  `json:"userId"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

var jwtAccessSecretKey = []byte("HelloIamJWTAccessSecretKey")
var jwtRefreshSecretKey = []byte("HelloIamJWTRefreshSecretKey")

func GenerateAuthTokens(userID int64, userRole string) (string, string, error) {
	op := "auth.GenerateAuthTokens"

	accessTokenPrepared := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenBody{
		userID,
		userRole,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	accessToken, err := accessTokenPrepared.SignedString(jwtAccessSecretKey)

	if err != nil {
		return "", "", fmt.Errorf("%v: error signed access token %v", op, err)
	}

	refreshTokenPrepared := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenBody{
		userID,
		userRole,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	refreshToken, err := refreshTokenPrepared.SignedString(jwtRefreshSecretKey)

	if err != nil {
		return "", "", fmt.Errorf("%v: error signed refresh token %v", op, err)
	}

	return accessToken, refreshToken, nil

}
