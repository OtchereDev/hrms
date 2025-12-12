package handlers

import (
	"github.com/OtchereDev/hrms-go/internal/api/middleware"
	"github.com/OtchereDev/hrms-go/internal/core/services/auth"
	"github.com/OtchereDev/hrms-go/pkg/response"
	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	authService *auth.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login handles user login
// POST /api/method/login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req auth.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	// Validate request
	if req.Email == "" || req.Password == "" {
		return response.BadRequest(c, "email and password are required")
	}

	// Get client IP and user agent
	ipAddress := c.IP()
	userAgent := c.Get("User-Agent")

	// Authenticate user
	loginResp, err := h.authService.Login(c.Context(), req, ipAddress, userAgent)
	if err != nil {
		switch err {
		case auth.ErrInvalidCredentials:
			return response.Unauthorized(c, "invalid email or password")
		case auth.ErrUserDisabled:
			return response.Forbidden(c, "user account is disabled")
		case auth.ErrUserLocked:
			return response.Forbidden(c, "user account is locked due to multiple failed login attempts")
		default:
			return response.InternalServerError(c, "failed to authenticate user")
		}
	}

	return response.Success(c, loginResp, "Logged In")
}

// Logout handles user logout
// POST /api/method/logout
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	sessionID := middleware.GetSessionID(c)
	if sessionID == "" {
		return response.BadRequest(c, "no active session")
	}

	err := h.authService.Logout(c.Context(), sessionID)
	if err != nil {
		return response.InternalServerError(c, "failed to logout")
	}

	return response.Success(c, nil, "Logged Out")
}

// RefreshToken handles token refresh
// POST /api/method/refresh_token
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if req.RefreshToken == "" {
		return response.BadRequest(c, "refresh token is required")
	}

	tokenPair, err := h.authService.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		switch err {
		case auth.ErrExpiredToken:
			return response.Unauthorized(c, "refresh token has expired")
		case auth.ErrInvalidToken:
			return response.Unauthorized(c, "invalid refresh token")
		case auth.ErrSessionExpired:
			return response.Unauthorized(c, "session has expired")
		case auth.ErrSessionRevoked:
			return response.Unauthorized(c, "session has been revoked")
		default:
			return response.InternalServerError(c, "failed to refresh token")
		}
	}

	return response.Success(c, tokenPair, "Token refreshed successfully")
}

// GetCurrentUserInfo returns current user information
// GET /api/method/hrms.api.get_current_user_info
func (h *AuthHandler) GetCurrentUserInfo(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Unauthorized(c, "not authenticated")
	}

	user, err := h.authService.GetUserByID(c.Context(), userID)
	if err != nil {
		return response.InternalServerError(c, "failed to get user info")
	}

	// Build roles list
	roles := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = role.Name
	}

	userInfo := map[string]interface{}{
		"id":         user.ID,
		"email":      user.Email,
		"username":   user.Username,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"full_name":  user.FullName,
		"mobile":     user.Mobile,
		"user_image": user.UserImage,
		"roles":      roles,
		"last_login": user.LastLogin,
		"enabled":    user.Enabled,
		"language":   user.Language,
		"time_zone":  user.TimeZone,
	}

	return response.Success(c, userInfo, "success")
}

// GetCurrentEmployeeInfo returns current employee information
// GET /api/method/hrms.api.get_current_employee_info
func (h *AuthHandler) GetCurrentEmployeeInfo(c *fiber.Ctx) error {
	// TODO: Implement when Employee model is created
	// For now, return empty response
	return response.Success(c, nil, "Employee info not yet implemented")
}
