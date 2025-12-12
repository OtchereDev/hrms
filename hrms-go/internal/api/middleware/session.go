package middleware

import (
	"strings"

	"github.com/OtchereDev/hrms-go/internal/core/services/session"
	"github.com/gofiber/fiber/v2"
)

// SessionMiddleware handles session-based authentication (for Frappe frontend)
type SessionMiddleware struct {
	sessionService *session.Service
}

// NewSessionMiddleware creates a new session middleware
func NewSessionMiddleware(sessionService *session.Service) *SessionMiddleware {
	return &SessionMiddleware{
		sessionService: sessionService,
	}
}

// Authenticate checks for valid session from cookies
func (m *SessionMiddleware) Authenticate(c *fiber.Ctx) error {
	// Get session ID from cookie
	sid := c.Cookies("sid")

	if sid == "" {
		// No session cookie, continue to next handler (might have JWT)
		return c.Next()
	}

	// Get session
	sess, err := m.sessionService.Get(c.Context(), sid)
	if err != nil {
		// Invalid session, clear cookie and continue
		c.ClearCookie("sid", "system_user", "full_name", "user_image")
		return c.Next()
	}

	// Store session info in context for handlers to use
	c.Locals("session", sess)
	c.Locals("user_id", sess.UserID)
	c.Locals("username", sess.User)
	c.Locals("email", sess.Email)
	c.Locals("full_name", sess.FullName)

	return c.Next()
}

// RequireSession requires a valid session (for session-only endpoints)
func (m *SessionMiddleware) RequireSession(c *fiber.Ctx) error {
	session := c.Locals("session")
	if session == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message":  nil,
			"exc":      "Not logged in",
			"exc_type": "AuthenticationError",
		})
	}

	return c.Next()
}

// DualAuth supports both JWT and session authentication
type DualAuthMiddleware struct {
	jwtMiddleware     *AuthMiddleware
	sessionMiddleware *SessionMiddleware
}

// NewDualAuthMiddleware creates middleware that accepts both JWT and sessions
func NewDualAuthMiddleware(jwtMiddleware *AuthMiddleware, sessionMiddleware *SessionMiddleware) *DualAuthMiddleware {
	return &DualAuthMiddleware{
		jwtMiddleware:     jwtMiddleware,
		sessionMiddleware: sessionMiddleware,
	}
}

// Authenticate tries JWT first, then session
func (m *DualAuthMiddleware) Authenticate(c *fiber.Ctx) error {
	// Check for JWT token in Authorization header
	authHeader := c.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		// Has JWT, use JWT auth
		return m.jwtMiddleware.Authenticate(c)
	}

	// No JWT, check for session cookie
	sid := c.Cookies("sid")
	if sid != "" {
		// Has session cookie, use session auth
		return m.sessionMiddleware.Authenticate(c)
	}

	// No authentication provided
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"message":  nil,
		"exc":      "Authentication required",
		"exc_type": "AuthenticationError",
	})
}
