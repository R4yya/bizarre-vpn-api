package handlers

import (
	"log/slog"
	"net/http"

	"bizarre-vpn-api/internal/api/helpers"
	"bizarre-vpn-api/internal/lib/logger/sl"
	"bizarre-vpn-api/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Log         *slog.Logger
	UserStorage services.UserStorage
}

type UserAuthorizationRequest struct {
	TelegramID   int64  `json:"telegramID" binding:"required"`
	Username     string `json:"username"`
	LanguageCode string `json:"languageCode"`
	IsBot        bool   `json:"isBot"`
}

// GetUserInfoHandler processes the user authorization request
// @Summary Get User Info
// @Security token
// @Description Getting base user data by token
// @Tags Users
// @Produce json
// @Success 200 {object} models.BaseUser "User Data"
// @Failure 401 {object} MessageResponse "Unauthorized"
// @Failure 500 {object} MessageResponse "Internal server error"
// @Router /users [get]
func (h *UserHandler) GetUserInfoHandler(c *gin.Context) {
	const op = "handlers.user.GetUserInfo"

	log := h.Log.With(
		slog.String("op", op),
	)

	tokenInfo, err := helpers.GetTokenInfo(c)

	if err != nil {
		log.Error("getting token info error", sl.Err(err))

		c.JSON(http.StatusBadRequest, MessageResponse{Message: err.Error()})
		return
	}

	userService := services.NewUserService(log, h.UserStorage)

	user, err := userService.GetUserById(tokenInfo.UserID)

	if err != nil {
		log.Error("getting user error", sl.Err(err))

		c.JSON(http.StatusBadRequest, MessageResponse{Message: err.Error()})
	}

	c.JSON(http.StatusOK, user)
}

// AuthorizeUserHandler processes the user authorization request
// @Summary User authorization
// @Description Authorize a user and register if it is not already in the database
// @Tags Users
// @Accept json
// @Produce json
// @Param user body UserAuthorizationRequest true "User Information"
// @Success 200 {object} models.User "The user authorized"
// @Success 201 {object} models.User "A new user has been successfully created"
// @Failure 400 {object} MessageResponse "Invalid request or missing required parameters"
// @Failure 409 {object} MessageResponse "User with this Telegram ID already exists"
// @Failure 500 {object} MessageResponse "Internal server error"
// @Router /users/auth [post]
// func (h *UserHandler) AuthorizeUserHandler(c *gin.Context) {
// 	const op = "handlers.user.AuthorizeUserHandler"

// 	log := h.Log.With(
// 		slog.String("op", op),
// 	)

// 	var req UserAuthorizationRequest

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, MessageResponse{Message: err.Error()})
// 		return
// 	}

// 	existingUser, err := h.userService.GetUser(req.TelegramID)
// 	if err == nil && existingUser != nil {
// 		c.JSON(http.StatusOK, existingUser)
// 		return
// 	}

// 	user := &models.User{
// 		TelegramID:   req.TelegramID,
// 		Username:     req.Username,
// 		LanguageCode: req.LanguageCode,
// 		IsBot:        req.IsBot,
// 	}

// 	userID, err := h.userService.RegisterUser(user)
// 	if err != nil {
// 		if errors.Is(err, storage.ErrUserAlreadyExists) {
// 			log.Info("user not found", sl.Err(err))
// 			c.JSON(http.StatusConflict, MessageResponse{Message: err.Error()})
// 			return
// 		}
// 		log.Error("registration user error", sl.Err(err))
// 		c.JSON(http.StatusInternalServerError, MessageResponse{Message: err.Error()})
// 		return
// 	}

// 	user.ID = userID

// 	c.JSON(http.StatusCreated, user)
// }
