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

// GetUserDataHandler processes the user authorization request
// @Summary Get User Data
// @Security token
// @scope.admin only administrative information
// @Description Getting base user list for admin
// @Tags Users
// @Produce json
// @Success 200 {object} []models.BaseUser "Users Data"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /users/list [get]
func (h *UserHandler) GetUsersListHandler(c *gin.Context) {
	const op = "handlers.user.GetUsersListHandler"

	log := h.Log.With(
		slog.String("op", op),
	)

	usersList, err := h.UserStorage.GetUsersList()

	if err != nil {
		log.Error("getting usersList error", sl.Err(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "something went wrong try again later"})
		return
	}

	c.JSON(http.StatusOK, usersList)
}

// GetUserDataHandler processes the user authorization request
// @Summary Get User Data
// @Security token
// @Description Getting base user data by token
// @Tags Users
// @Produce json
// @Success 200 {object} models.BaseUser "User Data"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /users [get]
func (h *UserHandler) GetUserDataHandler(c *gin.Context) {
	const op = "handlers.user.GetUserInfo"

	log := h.Log.With(
		slog.String("op", op),
	)

	tokenInfo, err := helpers.GetTokenInfo(c)

	if err != nil {
		log.Error("getting token info error", sl.Err(err))

		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	userService := services.NewUserService(log, h.UserStorage)

	user, err := userService.GetUserById(tokenInfo.UserID)

	if err != nil {
		log.Error("getting user error", sl.Err(err))

		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	c.JSON(http.StatusOK, user)
}
