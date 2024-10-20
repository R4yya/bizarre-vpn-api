package handlers

import (
	"bizarre-vpn-api/internal/storage/models"
	"bizarre-vpn-api/internal/storage/services"
	"bizarre-vpn-api/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserAuthorizationRequest struct {
	TelegramID   int64  `json:"telegram_id" binding:"required"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
	IsBot        bool   `json:"is_bot"`
}

type MessageResponse struct {
	Message string `json:"message"`
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
// @Failure 500 {object} MessageResponse "Internal server error"
// @Router /user/auth [post]
func AuthorizeUserHandler(c *gin.Context) {
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
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, MessageResponse{Message: err.Error()})
		return
	}

	user.ID = userID

	c.JSON(http.StatusCreated, user)
}
