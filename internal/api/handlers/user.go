package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"bizarre-vpn-api/internal/lib/logger/sl"
	"bizarre-vpn-api/internal/services"
	"bizarre-vpn-api/internal/storage"
	"bizarre-vpn-api/internal/storage/models"
)

type UserHandler struct {
	Log *slog.Logger
}

type UserAuthorizationRequest struct {
	TelegramID   int64  `json:"telegramID" binding:"required"`
	Username     string `json:"username"`
	LanguageCode string `json:"languageCode"`
	IsBot        bool   `json:"isBot"`
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
// @Router /user/auth [post]
func (h *UserHandler) AuthorizeUserHandler(c *gin.Context) {
	var req UserAuthorizationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, MessageResponse{Message: err.Error()})
		return
	}

	existingUser, err := services.GetUser(req.TelegramID)
	if err == nil && existingUser != nil {
		c.JSON(http.StatusOK, existingUser)
		return
	}

	user := &models.User{
		TelegramID:   req.TelegramID,
		Username:     req.Username,
		LanguageCode: req.LanguageCode,
		IsBot:        req.IsBot,
	}

	userID, err := services.RegisterUser(user)
	if err != nil {
		if errors.Is(err, storage.ErrUserAlreadyExists) {
			h.Log.Info("user not found", sl.Err(err))
			c.JSON(http.StatusConflict, MessageResponse{Message: err.Error()})
			return
		}
		h.Log.Error("registration user error", sl.Err(err))
		c.JSON(http.StatusInternalServerError, MessageResponse{Message: err.Error()})
		return
	}

	user.ID = userID

	c.JSON(http.StatusCreated, user)
}
