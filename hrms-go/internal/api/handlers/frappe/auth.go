package frappe

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	frappeCore "github.com/OtchereDev/hrms-go/internal/core/frappe"
	"github.com/OtchereDev/hrms-go/internal/core/services/auth"
	"github.com/OtchereDev/hrms-go/internal/core/services/session"
)

// AuthHandler handles Frappe-compatible authentication
type AuthHandler struct {
	authService    *auth.AuthService
	sessionService *session.Service
}

// NewAuthHandler creates a new Frappe auth handler
func NewAuthHandler(authService *auth.AuthService, sessionService *session.Service) *AuthHandler {
	return &AuthHandler{
		authService:    authService,
		sessionService: sessionService,
	}
}

// Login handles Frappe-style login
// POST /api/method/login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req auth.LoginRequest

	// Try JSON body first
	if err := c.BodyParser(&req); err != nil {
		// Try form data
		req.Email = c.FormValue("usr")
		req.Password = c.FormValue("pwd")
	}

	if req.Email == "" || req.Password == "" {
		return frappeCore.SendBadRequest(c, "username and password are required")
	}

	// Authenticate user
	ipAddress := c.IP()
	userAgent := c.Get("User-Agent")

	loginResp, err := h.authService.Login(c.Context(), req, ipAddress, userAgent)
	if err != nil {
		return frappeCore.SendAuthenticationError(c, "Invalid username or password")
	}

	// Create session
	sess, err := h.sessionService.Create(
		c.Context(),
		fmt.Sprintf("%d", loginResp.User.ID),
		loginResp.User.Username,
		loginResp.User.FullName,
		loginResp.User.Email,
		ipAddress,
		userAgent,
	)
	if err != nil {
		return frappeCore.SendInternalError(c, "failed to create session")
	}

	// Set cookies (Frappe compatible)
	c.Cookie(&fiber.Cookie{
		Name:     "sid",
		Value:    sess.SID,
		Path:     "/",
		HTTPOnly: true,
		SameSite: "Lax",
		Expires:  sess.ExpiresAt,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "system_user",
		Value:    sess.User,
		Path:     "/",
		HTTPOnly: false,
		SameSite: "Lax",
		Expires:  sess.ExpiresAt,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "full_name",
		Value:    sess.FullName,
		Path:     "/",
		HTTPOnly: false,
		SameSite: "Lax",
		Expires:  sess.ExpiresAt,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "user_id",
		Value:    sess.UserID,
		Path:     "/",
		HTTPOnly: false,
		SameSite: "Lax",
		Expires:  sess.ExpiresAt,
	})

	// Return Frappe-compatible response
	response := map[string]interface{}{
		"message":    "Logged In",
		"home_page":  "/app",
		"full_name":  sess.FullName,
		"user":       sess.User,
		"user_id":    sess.UserID,
		"email":      sess.Email,
		"session_id": sess.SID,
	}

	return frappeCore.SendSuccess(c, response)
}

// Logout handles Frappe-style logout
// POST /api/method/logout
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// Get session from cookie
	sid := c.Cookies("sid")

	if sid != "" {
		// Delete session
		_ = h.sessionService.Delete(c.Context(), sid)
	}

	// Clear cookies
	c.ClearCookie("sid", "system_user", "full_name", "user_id", "user_image")

	return frappeCore.SendSuccess(c, "Logged Out")
}

// GetLoggedUser returns the currently logged-in user
// GET /api/method/frappe.auth.get_logged_user
func (h *AuthHandler) GetLoggedUser(c *fiber.Ctx) error {
	// Check if user is logged in via session
	sess := c.Locals("session")
	if sess == nil {
		return frappeCore.SendAuthenticationError(c, "Not logged in")
	}

	session := sess.(*session.Session)

	response := map[string]interface{}{
		"user":      session.User,
		"user_id":   session.UserID,
		"full_name": session.FullName,
		"email":     session.Email,
	}

	return frappeCore.SendSuccess(c, response)
}

// GetSessionInfo returns session information
// GET /api/method/frappe.sessions.get_session_info
func (h *AuthHandler) GetSessionInfo(c *fiber.Ctx) error {
	// Check if user is logged in
	sess := c.Locals("session")
	if sess == nil {
		return frappeCore.SendSuccess(c, map[string]interface{}{
			"user":           "Guest",
			"session_exists": false,
		})
	}

	session := sess.(*session.Session)

	response := map[string]interface{}{
		"user":            session.User,
		"user_id":         session.UserID,
		"full_name":       session.FullName,
		"email":           session.Email,
		"session_exists":  true,
		"session_created": session.CreatedAt,
		"session_expires": session.ExpiresAt,
	}

	return frappeCore.SendSuccess(c, response)
}
