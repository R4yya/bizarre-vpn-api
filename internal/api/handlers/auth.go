package handlers

import (
	"bizarre-vpn-api/internal/api/helpers"
	"bizarre-vpn-api/internal/config"
	"bizarre-vpn-api/internal/services"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Log                    *slog.Logger
	CFG                    *config.Config
	UserStorage            services.UserStorage
	LnkUserProviderStorage services.LnkUserProviderStorage
}

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokensResponse struct {
	Message string `json:"message"`
	Tokens  Tokens `json:"tokens"`
}

// RefreshTokens processes refresh pair of auth tokens request
// @Summary Refresh Tokens
// @Description Generating new pair of auth tokens if refresh token is valid
// @Tags Users Auth
// @Accept json
// @Produce json
// @Param RefreshToken header string true "RefreshToken"
// @Success 200 {object} RefreshTokensResponse "Success generate new pair of tokens"
// @Failure 400 {object} MessageResponse "Invalid request or missing required parameters"
// @Failure 500 {object} MessageResponse "Internal server error"
// @Router /users/auth/refresh-tokens [post]
func (ah *AuthHandler) RefreshTokens(c *gin.Context) {
	const op = "handlers.auth.RefreshTokens"

	log := ah.Log.With(
		slog.String("op", op),
	)

	tokenInfo, token, err := helpers.ParseTokenFromHeader(c, log, "RefreshToken", ah.CFG.JWT.RefreshSecretKey)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	userService := services.NewUserService(log, ah.UserStorage)

	log.Debug("Getting user by userId from token", slog.Int64("userId", tokenInfo.UserID))
	user, err := userService.GetUserById(tokenInfo.UserID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	authService := services.NewAuthService(log, ah.UserStorage, ah.LnkUserProviderStorage)

	log.Debug("Getting refreshedTokens")
	accessToken, refreshToken, err := authService.GetRefreshedTokens(
		user,
		token,
		[]byte(ah.CFG.JWT.AccessSecretKey),
		[]byte(ah.CFG.JWT.RefreshSecretKey),
	)

	if err != nil {
		err := fmt.Errorf("getting user error: %w", err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, RefreshTokensResponse{
		Message: "Токены успешно обновлены",
		Tokens:  Tokens{AccessToken: accessToken, RefreshToken: refreshToken},
	})
}
