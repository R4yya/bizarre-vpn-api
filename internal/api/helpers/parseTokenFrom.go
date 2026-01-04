package helpers

import (
	"bizarre-vpn-api/internal/shared/jwt"
	"bizarre-vpn-api/internal/shared/logger/sl"
	"fmt"
	"log/slog"
	"strings"

	"github.com/gin-gonic/gin"
)

func ParseTokenFromHeader(
	c *gin.Context,
	log *slog.Logger,
	tokenKeyName string,
	jwtSecretKey string,
) (tokenInfo *jwt.TokenBody, receivedToken string, Err error) {
	authHeader := c.GetHeader(tokenKeyName)

	// Check if the header is present
	if authHeader == "" {
		return nil, "", fmt.Errorf("%v header is required", tokenKeyName)
	}

	log.Debug(fmt.Sprintf("Client %v header is received", tokenKeyName))

	return parseTokenFrom(log, tokenKeyName, authHeader, jwtSecretKey)
}

func ParseTokenFromCookie(
	c *gin.Context,
	log *slog.Logger,
	tokenKeyName string,
	jwtSecretKey string,
) (tokenInfo *jwt.TokenBody, receivedToken string, Err error) {
	cookieString, err := c.Cookie(tokenKeyName)

	if err != nil {
		return nil, "", fmt.Errorf("%v cookie is required: %v", tokenKeyName, err)
	}

	if cookieString == "" {
		return nil, "", fmt.Errorf("%v cookie is required", tokenKeyName)
	}

	log.Debug(fmt.Sprintf("Client %v cookie is received", tokenKeyName))

	return parseTokenFrom(log, tokenKeyName, cookieString, jwtSecretKey)
}

func parseTokenFrom(
	log *slog.Logger,
	tokenKeyName string,
	tokenString string,
	jwtSecretKey string,
) (tokenInfo *jwt.TokenBody, receivedToken string, Err error) {
	parts := strings.Split(tokenString, " ")

	if parts[0] != "Bearer" {
		return nil, "", fmt.Errorf("incorrect flow of %v token, support only a bearer flow", tokenKeyName)
	}

	log.Debug("Client bearer flow is correctly")

	token := parts[1]

	if token == "" {
		return nil, "", fmt.Errorf("%v token must be a noempty string", tokenKeyName)
	}

	log.Debug("token string is noempty")

	tokenInfo, err := jwt.ParseToken(token, []byte(jwtSecretKey))

	log.Debug("parseFromToken", slog.Any("tokenInfo", tokenInfo))

	if err != nil {
		log.Error(fmt.Sprintf("%v token parse error", tokenKeyName), sl.Err(err))

		return nil, "", fmt.Errorf("%v token is invalid", tokenKeyName)
	}

	return tokenInfo, token, nil
}
