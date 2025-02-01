package handlers

import (
	"log/slog"

	"bizarre-vpn-api/internal/services"
)

type UserHandler struct {
	Log         *slog.Logger
	userService *services.UserService
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
