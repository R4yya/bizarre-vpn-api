package handlers

import (
	"bizarre-vpn-api/internal/api/helpers"
	"bizarre-vpn-api/internal/config"
	"bizarre-vpn-api/internal/lib/initData"
	"bizarre-vpn-api/internal/lib/logger/sl"
	"bizarre-vpn-api/internal/services"
	"fmt"
	"log/slog"
	"net/http"
	"time"

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

type TokensResponse struct {
	Message string `json:"message"`
	Tokens  Tokens `json:"tokens"`
}

type InitDataRequestData struct {
	InitDataStr string `json:"initDataStr"`
}

// AuthorizeWithInitData processes the user authorization with telegram initData
// @Summary User authorization
// @Description Authorize a user and register if it is not already in the database
// @Tags Users Auth
// @Accept json
// @Produce json
// @Param initDataStr body InitDataRequestData true "telegram user initData string"
// @Success 200 {object} TokensResponse "Success generate new pair of tokens"
// @Failure 400 {object} MessageResponse "Invalid request or missing required parameters"
// @Failure 500 {object} MessageResponse "Internal server error"
// @Router /users/auth/telegram-init-data [post]
func (ah *AuthHandler) AuthorizeWithInitData(c *gin.Context) {
	const op = "handlers.auth.AuthorizeWithInitData"

	log := ah.Log.With(
		slog.String("op", op),
	)

	var req InitDataRequestData

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, MessageResponse{Message: err.Error()})
		return
	}

	initDataPayload, isOk, err := initData.ValidateInitData(req.InitDataStr, ah.CFG.TelegramBotToken, time.Hour)

	if err != nil || !isOk {
		log.Info("ValidateInitData err", sl.Err(err))
		c.JSON(http.StatusBadRequest, MessageResponse{Message: "ValidateInitData error"})
		return
	}

	authService := services.NewAuthService(
		log,
		ah.UserStorage,
		ah.LnkUserProviderStorage,
	)

	accessToken, refreshToken, err := authService.AuthorizeByTelegram(
		initDataPayload.TelegramID,
		initDataPayload.Username,
		[]byte(ah.CFG.JWT.AccessSecretKey),
		[]byte(ah.CFG.JWT.RefreshSecretKey),
	)

	if err != nil {
		log.Error(fmt.Sprintf("%v: %v", op, "AuthorizeByTelegram"), sl.Err(err))
		c.JSON(http.StatusInternalServerError, MessageResponse{Message: "error: internal server error, try again later"})
		return
	}

	c.JSON(http.StatusOK, TokensResponse{
		Message: "Successful authorize with telegram",
		Tokens:  Tokens{AccessToken: accessToken, RefreshToken: refreshToken},
	})
}

// RefreshTokens processes refresh pair of auth tokens request
// @Summary Refresh Tokens
// @Description Generating new pair of auth tokens if refresh token is valid
// @Tags Users Auth
// @Accept json
// @Produce json
// @Param RefreshToken header string true "RefreshToken"
// @Success 200 {object} TokensResponse "Success generate new pair of tokens"
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

	c.JSON(http.StatusOK, TokensResponse{
		Message: "Токены успешно обновлены",
		Tokens:  Tokens{AccessToken: accessToken, RefreshToken: refreshToken},
	})
}
