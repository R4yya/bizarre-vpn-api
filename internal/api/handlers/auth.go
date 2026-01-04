package handlers

import (
	"bizarre-vpn-api/internal/api/helpers"
	"bizarre-vpn-api/internal/services"
	"bizarre-vpn-api/internal/services/interfaces"
	"bizarre-vpn-api/internal/shared/config"
	"bizarre-vpn-api/internal/shared/coreErrors"
	"bizarre-vpn-api/internal/shared/logger/sl"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Log                    *slog.Logger
	CFG                    *config.Config
	UserStorage            interfaces.UserStorage
	LnkUserProviderStorage interfaces.LnkUserProviderStorage
}

type AuthResponse struct {
	Message     string `json:"message"`
	AccessToken string `json:"accessToken"`
}

type InitDataRequestData struct {
	InitDataStr string `json:"initDataStr"`
}

type CredentialRequestData struct {
	Login    string
	Password string
}

const RefreshTokenCookieKey = "refreshToken"

var maxAge = 60 * 60 * 24 * 60 // lifetime in seconds
var expirationTime = time.Now().Add(time.Duration(maxAge) * time.Second).UTC().Format(http.TimeFormat)

func setNewRefreshToken(c *gin.Context, refreshToken string) {
	tokenString := "Bearer " + refreshToken

	c.Writer.Header().Set("Set-Cookie", fmt.Sprintf("%v=%v; Path=/; Domain=; Secure; HttpOnly; SameSite=None; Max-Age=%v; Expires=%v", RefreshTokenCookieKey, tokenString, maxAge, expirationTime))
}

// // AuthorizeWithInitData processes the user authorization with telegram initData
// // @Summary User authorization
// // @Description Authorize a user and register if it is not already in the database
// // @Tags Users Auth
// // @Accept json
// // @Produce json
// // @Param initDataStr body InitDataRequestData true "telegram user initData string"
// // @Success 200 {object} AuthResponse "Success generate new pair of tokens"
// // @Failure 400 {object} ErrorResponse "Invalid request or missing required parameters"
// // @Failure 500 {object} ErrorResponse "Internal server error"
// // @Router /users/auth/telegram-init-data [post]
// func (ah *AuthHandler) AuthorizeWithInitData(c *gin.Context) {
// 	const op = "handlers.auth.AuthorizeWithInitData"

// 	log := ah.Log.With(
// 		slog.String("op", op),
// 	)

// 	var req InitDataRequestData

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
// 		return
// 	}

// 	initDataPayload, err := initData.ValidateInitData(req.InitDataStr, ah.CFG.TelegramBot.BotToken, time.Hour)

// 	if err != nil {
// 		log.Info("ValidateInitData err", sl.Err(err))
// 		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "ValidateInitData error"})
// 		return
// 	}

// 	authService := services.NewAuthService(
// 		log,
// 		ah.UserStorage,
// 		ah.LnkUserProviderStorage,
// 	)

// 	accessToken, refreshToken, err := authService.AuthorizeByTelegram(
// 		initDataPayload.TelegramID,
// 		initDataPayload.Username,
// 		[]byte(ah.CFG.JWT.AccessSecretKey),
// 		[]byte(ah.CFG.JWT.RefreshSecretKey),
// 	)

// 	if err != nil {
// 		log.Error(fmt.Sprintf("%v: %v", op, "AuthorizeByTelegram"), sl.Err(err))
// 		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error, try again later"})
// 		return
// 	}

// 	setNewRefreshToken(c, refreshToken)

// 	c.JSON(http.StatusOK, AuthResponse{
// 		Message:     "Successful authorize with telegram",
// 		AccessToken: accessToken,
// 	})
// }

// AuthorizeWithCredentials processes the user authorization with username and password
// @Summary User authorization
// @Description User authorization with username and password
// @Tags Users Auth
// @Accept json
// @Produce json
// @Param credentials body CredentialRequestData true "Username and Password strings"
// @Success 200 {object} AuthResponse "Success generate new pair of tokens"
// @Failure 400 {object} ErrorResponse "Incorrect username or password"
// @Failure 400 {object} ErrorResponse "Invalid request or missing required parameters"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /users/auth/credentials [post]
func (ah *AuthHandler) AuthorizeWithCredentials(c *gin.Context) {
	const op = "handlers.auth.AuthorizeWithInitData"

	log := ah.Log.With(
		slog.String("op", op),
	)

	var req CredentialRequestData

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	authService := services.NewAuthService(
		log,
		ah.UserStorage,
		ah.LnkUserProviderStorage,
	)

	accessToken, refreshToken, err := authService.AuthorizeByCredentials(
		req.Login,
		req.Password,
		[]byte(ah.CFG.JWT.AccessSecretKey),
		[]byte(ah.CFG.JWT.RefreshSecretKey),
	)

	if err != nil {
		if errors.Is(err, coreErrors.ErrorIncorrectLoginOrPass) {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			c.Abort()
			return
		}

		log.Error("authorize error", sl.Err(err))

		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Something went wrong"})
		c.Abort()
		return
	}

	setNewRefreshToken(c, refreshToken)

	c.JSON(http.StatusOK, AuthResponse{
		Message:     "Successful authorize with credentials",
		AccessToken: accessToken,
	})
}

// RefreshTokens processes refresh pair of auth tokens request
// @Summary Refresh Tokens
// @Description Generating new pair of auth tokens if refresh token is valid
// @Tags Users Auth
// @Produce json
// @Success 200 {object} AuthResponse "Success generate new pair of tokens"
// @Failure 400 {object} ErrorResponse "Invalid request or missing required parameters"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /users/auth/refresh-tokens [post]
func (ah *AuthHandler) RefreshTokens(c *gin.Context) {
	const op = "handlers.auth.RefreshTokens"

	log := ah.Log.With(
		slog.String("op", op),
	)

	tokenInfo, token, err := helpers.ParseTokenFromCookie(c, log, RefreshTokenCookieKey, ah.CFG.JWT.RefreshSecretKey)

	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		c.Abort()
		return
	}

	userService := services.NewUserService(log, ah.UserStorage)

	log.Debug("Getting user by userId from token", slog.Int64("userId", tokenInfo.UserID))
	user, err := userService.GetUserById(tokenInfo.UserID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
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

		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		c.Abort()
		return
	}

	setNewRefreshToken(c, refreshToken)

	c.JSON(http.StatusOK, AuthResponse{
		Message:     "Токены успешно обновлены",
		AccessToken: accessToken,
	})
}
