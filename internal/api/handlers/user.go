package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"bizarre-vpn-api/internal/api/helpers"
	"bizarre-vpn-api/internal/bot"
	"bizarre-vpn-api/internal/lib/logger/sl"
	"bizarre-vpn-api/internal/services/authLinkService"
	"bizarre-vpn-api/internal/services/userService"
	"bizarre-vpn-api/internal/storage/models"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Log             *slog.Logger
	BotSharedData   *bot.BotSharedData
	UserService     *userService.UserService
	AuthLinkService *authLinkService.AuthLinkService
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

	usersList, err := h.UserService.GetUsersList()

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

	user, err := h.UserService.GetUserById(tokenInfo.UserID)

	if err != nil {
		log.Error("getting user error", sl.Err(err))

		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	c.JSON(http.StatusOK, user)
}

// CreateUser processes the user authorization request
// @Summary Create User
// @Security token
// @scope.admin only administrative information
// @Description Create User by admin
// @Tags Users
// @Accept json
// @Produce json
// @Param CreateUserPayload body models.CreateUserPayload true "Create User Data"
// @Success 200 {object} []models.BaseUser "Users Data"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	const op = "handlers.user.CreateUser"

	log := h.Log.With(
		slog.String("op", op),
	)

	tokenInfo, err := helpers.GetTokenInfo(c)

	if err != nil {
		log.Error("getting token info error", sl.Err(err))

		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		c.Abort()
		return
	}

	if tokenInfo.Role != models.UserRoleAdmin {
		mes := "create user forbidder for no admin users"
		log.Debug(mes)

		c.JSON(http.StatusBadRequest, ErrorResponse{Error: mes})
		c.Abort()
		return
	}

	var body models.CreateUserPayload

	err = c.ShouldBindJSON(&body)

	if err != nil {
		log.Error("parsing request json error", sl.Err(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		c.Abort()
		return
	}

	log.Debug("body", slog.Any("CreateUserPayload", body))

	createdUser, err := h.UserService.CreateUser(&body)

	if err != nil {
		if errors.Is(userService.ErrIncorrectRole, err) ||
			errors.Is(userService.ErrInvalidPassword, err) ||
			errors.Is(userService.ErrLoginOccupied, err) ||
			errors.Is(userService.ErrLoginIsTooSmall, err) {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			c.Abort()
			return
		}

		log.Error("create user error", sl.Err(err))

		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "something went wrong"})
		c.Abort()

		return
	}

	c.JSON(http.StatusOK, createdUser)
}

type UrlResponse struct {
	Url string `json:"url"`
}

// CreateUserAuthLink returns tg bot auth url
// @Summary CreateUserAuthLink by userId
// @Description CreateUserAuthLink by userId
// @Security token
// @scope.admin only administrative information
// @Tags Users
// @Param id path int true "UserId"
// @Produce json
// @Success 200 {object} UrlResponse "tg bot auth url"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 400 {object} ErrorResponse "User already linked"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /users/{id}/auth-link [get]
func (h *UserHandler) CreateUserAuthLink(c *gin.Context) {
	const op = "handlers.user.CreateUserAuthLink"

	log := h.Log.With(
		slog.String("op", op),
	)

	userIdStr := c.Param("id")

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		err = fmt.Errorf("convert error vpnServerIdStr: %w", err)
		log.Info(err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		c.Abort()
		return
	}

	authLink, err := h.AuthLinkService.CreateItem(userId)

	if err != nil {
		if errors.Is(authLinkService.ErrUserAlreadyLinked, err) {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			c.Abort()
			return
		}

		log.Error("create user AuthLink error", sl.Err(err))

		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "something went wrong"})
		c.Abort()

		return
	}

	url := fmt.Sprintf("https://t.me/%v?start=%v", h.BotSharedData.Username, authLink.Code)

	c.JSON(http.StatusOK, UrlResponse{Url: url})
}
